package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/blockchain-tps-test/samples/theta/tps"
	"github.com/pkg/errors"
	"github.com/thetatoken/theta/blockchain"
	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/crypto"
	"github.com/thetatoken/theta/ledger/types"
	"github.com/thetatoken/theta/rpc"
	"github.com/ybbus/jsonrpc"

	ct "github.com/blockchain-tps-test/samples/theta/accessors"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/core"
	"github.com/thetatoken/thetasubchain/eth/abi"
	"github.com/thetatoken/thetasubchain/eth/abi/bind"
	"github.com/thetatoken/thetasubchain/eth/ethclient"
	rpcc "github.com/ybbus/jsonrpc"
)

var (
	_ tps.Client = (*EthClient)(nil)
)
var dec18, _ = new(big.Int).SetString("1000000000000000000", 10)
var crossChainFee = new(big.Int).Mul(big.NewInt(10), dec18)

type EthClient struct {
	rpcClient *jsonrpc.RPCClient //theta-rpc
	client    *ethclient.Client  //eth-adaptor

}

func NewClient(rpcClientUrl ,ethClientUrl string) (c EthClient, err error) { //,theta-url,eth-url
	c.client, err = ethclient.Dial(ethClientUrl)
	if err != nil {
		return
	}
	c.rpcClient = jsonrpc.NewRPCClient(rpcClientUrl)
	return
}

type GetStatusArgs struct{}

type GetStatusResult struct {
	Address                    string            `json:"address"`
	ChainID                    string            `json:"chain_id"`
	PeerID                     string            `json:"peer_id"`
	LatestFinalizedBlockHash   common.Hash       `json:"latest_finalized_block_hash"`
	LatestFinalizedBlockHeight common.JSONUint64 `json:"latest_finalized_block_height"`
	LatestFinalizedBlockTime   *common.JSONBig   `json:"latest_finalized_block_time"`
	LatestFinalizedBlockEpoch  common.JSONUint64 `json:"latest_finalized_block_epoch"`
	CurrentEpoch               common.JSONUint64 `json:"current_epoch"`
	CurrentHeight              common.JSONUint64 `json:"current_height"`
	CurrentTime                *common.JSONBig   `json:"current_time"`
	Syncing                    bool              `json:"syncing"`
	GenesisBlockHash           common.Hash       `json:"genesis_block_hash"`
	SnapshotBlockHeight        common.JSONUint64 `json:"snapshot_block_height"`
	SnapshotBlockHash          common.Hash       `json:"snapshot_block_hash"`
}

func HandleThetaRPCResponse(rpcRes *rpcc.RPCResponse, rpcErr error, parse func(jsonBytes []byte) (interface{}, error)) (result interface{}, err error) {
	if rpcErr != nil {
		return nil, fmt.Errorf("failed to get theta RPC response: %v", rpcErr)
	}
	if rpcRes.Error != nil {
		return nil, fmt.Errorf("theta RPC returns an error: %v", rpcRes.Error)
	}

	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(rpcRes.Result, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to parse theta RPC response: %v, %s", err, string(jsonBytes))
	}

	//logger.Infof("HandleThetaRPCResponse, jsonBytes: %v", string(jsonBytes))
	result, err = parse(jsonBytes)
	if err != nil {
		logger.Warn("Failed to parse theta RPC response: " + err.Error())
	}
	return
}
func (c EthClient) LatestBlockHeight(ctx context.Context) (uint64, error) {

	// res, err := c.client.BlockNumber(ctx)
	// if err != nil {
	// 	return 0, err
	// }
	// return res, nil
	rpcRes, rpcErr := c.rpcClient.Call("theta.GetStatus", GetStatusArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.LatestFinalizedBlockHeight, nil
	}
	var height tcommon.JSONUint64
	resultIntf, err := HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	height = resultIntf.(tcommon.JSONUint64)
	result := uint64(height)
	if err != nil {
		return result, err
	}
	return result, nil

}

type EthGetBlockResult struct {
	Height    hexutil.Uint64 `json:"number"`
	Hash      common.Hash    `json:"hash"`
	Parent    common.Hash    `json:"parentHash"`
	Timestamp hexutil.Uint64 `json:"timestamp"`
	Proposer  common.Address `json:"miner"`
	TxHash    common.Hash    `json:"transactionsRoot"`
	StateHash common.Hash    `json:"stateRoot"`

	ReiceptHash     common.Hash    `json:"receiptsRoot"`
	Nonce           string         `json:"nonce"`
	Sha3Uncles      common.Hash    `json:"sha3Uncles"`
	LogsBloom       string         `json:"logsBloom"`
	Difficulty      hexutil.Uint64 `json:"difficulty"`
	TotalDifficulty hexutil.Uint64 `json:"totalDifficulty"`
	Size            hexutil.Uint64 `json:"size"`
	GasLimit        hexutil.Uint64 `json:"gasLimit"`
	GasUsed         hexutil.Uint64 `json:"gasUsed"`
	ExtraData       string         `json:"extraData"`
	Uncles          []common.Hash  `json:"uncles"`
	Transactions    []interface{}  `json:"transactions"`
}
type ThetaGetBlockResult struct {
	*ThetaGetBlockResultInner
}
type ThetaGetBlocksResult []*ThetaGetBlockResultInner

type ThetaGetBlockResultInner struct {
	ChainID            string                   `json:"chain_id"`
	Epoch              common.JSONUint64        `json:"epoch"`
	Height             common.JSONUint64        `json:"height"`
	Parent             common.Hash              `json:"parent"`
	TxHash             common.Hash              `json:"transactions_hash"`
	StateHash          common.Hash              `json:"state_hash"`
	Timestamp          *common.JSONBig          `json:"timestamp"`
	Proposer           common.Address           `json:"proposer"`
	HCC                core.CommitCertificate   `json:"hcc"`
	GuardianVotes      *core.AggregatedVotes    `json:"guardian_votes"`
	EliteEdgeNodeVotes *core.AggregatedEENVotes `json:"elite_edge_node_votes"`

	Children []common.Hash    `json:"children"`
	Status   core.BlockStatus `json:"status"`

	Hash common.Hash `json:"hash"`
	Txs  []rpc.Tx    `json:"transactions"`
}
type LogData struct {
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`
}
type RPCResult struct {
	Result  []LogData `json:"logs"`
	Address string    `json:"contractAddress"`
}

func parse(jsonBytes []byte) (int, time.Duration, error) {
	trpcResult := ThetaGetBlockResult{}
	json.Unmarshal(jsonBytes, &trpcResult)
	result := 0
	if trpcResult.ThetaGetBlockResultInner == nil {
		return result, time.Duration(0), errors.New("empty block")
	}

	var objmap map[string]json.RawMessage
	var elapsedTime time.Duration
	json.Unmarshal(jsonBytes, &objmap)
	if objmap["transactions"] != nil {
		var txmaps []map[string]json.RawMessage
		json.Unmarshal(objmap["transactions"], &txmaps)
		for i, value := range txmaps {
			if types.TxType(trpcResult.Txs[i].Type) == types.TxSmartContract {
				mutex.Lock()
				startTime := txMapCrossChain[countChainTx2.String()]
				var test RPCResult
				json.Unmarshal(value["receipt"], &test)
				fmt.Println(test.Result[1].Topics)
				fmt.Println(crypto.Keccak256Hash([]byte("TNT20VoucherMinted(string,address,address,uint256,uint256,uint256)")).Hex())
				type TransferEvt struct {
					Denom                      string
					TargetChainVoucherReceiver common.Address
					VoucherContact             common.Address
					MintedAmount               *big.Int
					SourceChainTokenLockNonce  *big.Int
					VoucherMintNonce           *big.Int
				}
				var event TransferEvt
				contractAbi, _ := abi.JSON(strings.NewReader(RawABI))
				txData := test.Result[1].Data
				h := []byte(txData)

				err := contractAbi.UnpackIntoInterface(&event, "TNT20VoucherMinted", h)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(event.VoucherMintNonce)
				countChainTx2.Add(countChainTx2, big.NewInt(1))
				mutex.Unlock()
				elapsedTime = elapsedTime + time.Since(startTime)/time.Second
				result += 1
			} else if types.TxType(trpcResult.Txs[i].Type) == types.TxSend {
				mutex.Lock()
				startTime := txMap[trpcResult.Txs[i].Hash]
				mutex.Unlock()
				elapsedTime = elapsedTime + time.Since(startTime)/time.Millisecond
				//fmt.Println(elapsedTime)
				result += 1
			}
		}
	}
	if result != 0 {
		avgLatency = elapsedTime / time.Duration(result)
	}
	//fmt.Println(avgLatency)
	return result, avgLatency, nil
}
func (c EthClient) CountTx(ctx context.Context, height uint64) (int, time.Duration, error) {
	startTime := time.Now()
	rpcResult, err := c.rpcClient.Call("theta.GetBlockByHeight", rpc.GetBlockByHeightArgs{
		Height: common.JSONUint64(height)})
	if err != nil {
		return 0, time.Duration(0), err
	}
	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(rpcResult.Result, "", "    ")

	//logger.Infof("HandleThetaRPCResponse, jsonBytes: %v", strin(jsonBytes))
	result, avg_latency, err := parse(jsonBytes)
	totalTime := time.Since(startTime) / time.Millisecond
	fmt.Println("call and parse consume ", totalTime)
	return result, avg_latency, nil
}

func (c EthClient) CountPendingTx(ctx context.Context) (int, error) {

	// count, err := c.client.PendingTransactionCount(ctx)
	// if err != nil {
	// 	return 0, err
	// }
	return int(0), nil
}

type GetAccountArgs struct {
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Height  common.JSONUint64 `json:"height"`
	Preview bool              `json:"preview"` // preview the account balance from the ScreenedView
}
type GetAccountResult struct {
	*types.Account
	Address string `json:"address"`
}

func (c EthClient) Nonce(ctx context.Context, address string) (uint64, error) {

	// nonce, err := c.client.PendingNonceAt(ctx, common.HexToAddress(address))
	// if err != nil {
	// 	log.Fatalln("error in getNonce while calling pending nonce at", err)

	// 	return 0, err
	// }
	// return nonce, nil
	//height := tcommon.JSONUint64(math.MaxUint64)
	rpcRes, rpcErr := c.rpcClient.Call("theta.GetAccount", GetAccountArgs{Address: address, Preview: true})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := GetAccountResult{Account: &types.Account{}}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Account.Sequence, nil
	}

	resultIntf, err := HandleThetaRPCResponse(rpcRes, rpcErr, parse)

	if err != nil {
		return 0, nil
	}

	// result = fmt.Sprintf("0x%x", resultIntf.(*big.Int))
	result := resultIntf.(uint64)
	return result, nil
}

type chainIDResultWrapper struct {
	chainID string
}

func (c *EthClient) getChainID(ctx context.Context) big.Int {

	// chainid, err := c.client.ChainID(ctx)
	// if err != nil {
	// 	log.Fatalln("error in getChainID while getting chain id:", err)
	// 	return big.Int{}
	// }
	// return *chainid
	rpcRes, rpcErr := c.rpcClient.Call("theta.GetStatus", GetStatusArgs{})
	var blockHeight uint64
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := chainIDResultWrapper{
			chainID: trpcResult.ChainID,
		}
		blockHeight = uint64(trpcResult.LatestFinalizedBlockHeight)
		return re, nil
	}

	resultIntf, _ := HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	thetaChainIDResult, _ := resultIntf.(chainIDResultWrapper)

	thetaChainID := thetaChainIDResult.chainID
	ethChainID := types.MapChainID(thetaChainID, blockHeight).Uint64()
	return *big.NewInt(int64(ethChainID))
}

type GetBlockByHeightArgs struct {
	Height             common.JSONUint64 `json:"height"`
	IncludeEthTxHashes bool              `json:"include_eth_tx_hashes"`
}
type Tx struct {
	types.Tx       `json:"raw"`
	Type           byte                              `json:"type"`
	Hash           common.Hash                       `json:"hash"`
	Receipt        *blockchain.TxReceiptEntry        `json:"receipt"`
	BalanceChanges *blockchain.TxBalanceChangesEntry `json:"balance_changes"`
}
type TxTmp struct {
	Tx   json.RawMessage `json:"raw"`
	Type byte            `json:"type"`
	Hash tcommon.Hash    `json:"hash"`
}

func (c *EthClient) getGasPriceSuggestion(ctx context.Context) big.Int {

	// gasPrice, err := c.client.SuggestGasPrice(ctx)
	// if err != nil {
	// 	log.Fatalln("error in SuggestGasPrice while calling pending nonce at", err)

	// 	return big.Int{}
	// }
	// return *gasPrice
	currentHeight, _ := c.LatestBlockHeight(ctx)
	rpcRes, rpcErr := c.rpcClient.Call("theta.GetBlockByHeight", GetBlockByHeightArgs{Height: tcommon.JSONUint64(currentHeight)})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := ThetaGetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		var objmap map[string]json.RawMessage
		json.Unmarshal(jsonBytes, &objmap)
		if objmap["transactions"] != nil {
			//TODO: handle other types
			txs := []rpc.Tx{}
			tmpTxs := []TxTmp{}
			json.Unmarshal(objmap["transactions"], &tmpTxs)
			for _, tx := range tmpTxs {
				newTx := rpc.Tx{}
				newTx.Type = tx.Type
				newTx.Hash = tx.Hash
				if types.TxType(tx.Type) == types.TxSmartContract {
					transaction := types.SmartContractTx{}
					json.Unmarshal(tx.Tx, &transaction)
					// fmt.Printf("transaction: %+v\n", transaction)
					newTx.Tx = &transaction
				}
				txs = append(txs, newTx)
			}
			trpcResult.Txs = txs
		}
		return trpcResult, nil
	}

	resultIntf, _ := HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	thetaGetBlockResult, _ := resultIntf.(ThetaGetBlockResult)

	totalGasPrice := big.NewInt(0)
	count := 0
	for _, tx := range thetaGetBlockResult.Txs {
		if types.TxType(tx.Type) != types.TxSmartContract {
			continue
		}
		if tx.Tx != nil {
			transaction := tx.Tx.(*types.SmartContractTx)
			count++
			totalGasPrice = new(big.Int).Add(transaction.GasPrice, totalGasPrice)
		}
	}
	gasPrice := big.NewInt(4000000000000)
	if count != 0 {
		gasPrice = new(big.Int).Div(totalGasPrice, big.NewInt(int64(count)))
	}
	//.Printf("gasPrice: %v\n", gasPrice)
	//result := "0x" + gasPrice.Text(16)
	return *gasPrice

}

// sends transaction to the network
func (c *EthClient) SendTx(ctx context.Context, privHex string, nonce uint64, to string, value int64) (common.Hash, error) {

	// signedtx, err := c.signTransaciton(ctx, privHex, nonce, to, value)

	// if err != nil {
	// 	return common.BytesToHash([]byte("")), err
	// }
	// if err := c.client.SendTransaction(ctx, signedtx); err != nil {
	// 	log.Fatalln("error in SendTx while getting chain id:", err)

	// }
	// fmt.Println("transaction sent. txid: ", signedtx.Hash().Hex(), "nonce: ", nonce)
	//wallet, address, err := tx.SoftWalletUnlock("/home/dd/.thetacli", "2E833968E5bB786Ae419c4d13189fB081Cc43bab", "qwertyuiop")
	time.Sleep(1 * time.Millisecond)
	privateKey, err := crypto.HexToECDSA(privHex)
	pri, err := hex.DecodeString(privHex)
	thetaPrivateKey, err := crypto.PrivateKeyFromBytes(pri)
	// privateKey, err := crypto.HexToECDSA("2dad160420b1e9b6fc152cd691a686a7080a0cee41b98754597a2ce57cc5dab1")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := pubkeyToAddress(*publicKeyECDSA)
	//nonce, _ = c.Nonce(context.Background(), fromAddress.String())
	nonce -= 1
	theta := big.NewInt(0)
	tfuel := big.NewInt(1)
	fee := big.NewInt(3e17)
	inputs := []types.TxInput{{
		Address: fromAddress,
		Coins: types.Coins{
			TFuelWei: new(big.Int).Add(tfuel, fee),
			ThetaWei: theta,
		},
		Sequence: uint64(nonce + 1),
	}}
	//fmt.Println("receive", nonce+1)
	outputs := []types.TxOutput{{
		Address: common.HexToAddress("19E7E376E7C213B7E7e7e46cc70A5dD086DAff3A"),
		Coins: types.Coins{
			TFuelWei: tfuel,
			ThetaWei: theta,
		},
	}}
	sendTx := &types.SendTx{
		Fee: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			TFuelWei: fee,
		},
		Inputs:  inputs,
		Outputs: outputs,
	}
	sig, err := thetaPrivateKey.Sign(sendTx.SignBytes("testnet")) //privatenet,testnet
	if err != nil {
		log.Fatalln("Failed to sign transaction: %v\n", err)
	}
	sendTx.SetSignature(fromAddress, sig)

	raw, err := types.TxToBytes(sendTx)
	if err != nil {
		log.Fatalln("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)
	var res *jsonrpc.RPCResponse
	res, err = c.rpcClient.Call("theta.BroadcastRawTransactionAsync", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	if err != nil {
		return common.BytesToHash(nil), err
		log.Fatalln("Failed to broadcast transaction: %v\n", err)
	}
	if res.Error != nil {
		return common.BytesToHash(nil), res.Error
		log.Fatalln("Server returned error: %v\n", res.Error)
	}
	result := &rpc.BroadcastRawTransactionResult{}
	err = res.GetObject(result)
	if err != nil {
		return common.BytesToHash(nil), err
		log.Fatalln("Failed to parse server response: %v\n", err)
	}
	formatted, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return common.BytesToHash(formatted), err
		log.Fatalln("Failed to parse server response: %v\n", err)
	}
	//fmt.Printf("Successfully broadcasted transaction:\n%s\n", formatted)
	startTime := time.Now()
	mutex.Lock()
	txMap[common.HexToHash(result.TxHash)] = startTime
	mutex.Unlock()
	CountNum += 1
	// if CountNum%100 == 0 {
	// 	fmt.Println("already send ", CountNum)
	// }
	fmt.Println("have send ", CountNum, " txs", nonce+1)
	return common.Hash{}, err //common.BytesToHash(formatted), err
}

const RawABI = `
[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "denom",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "targetChainVoucherReceiver",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "voucherContact",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "mintedAmount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "sourceChainTokenLockNonce",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "voucherMintNonce",
				"type": "uint256"
			}
		],
		"name": "TNT20VoucherMinted",
		"type": "event"
	}
]`

func Resolve(data []byte) big.Int {
	contractAbi, err := abi.JSON(strings.NewReader(RawABI))
	if err != nil {
		fmt.Println(err)
	}
	type TransferEvt struct {
		Denom                      string
		TargetChainVoucherReceiver common.Address
		VoucherContact             common.Address
		MintedAmount               *big.Int
		SourceChainTokenLockNonce  *big.Int
		VoucherMintNonce           *big.Int
	}
	var event TransferEvt
	h := data
	err = contractAbi.UnpackIntoInterface(&event, "TNT20VoucherMinted", h)
	if err != nil {
		fmt.Println(err)
	}
	return *event.VoucherMintNonce
}
func (c EthClient) Erc20TransferFrom(ctx context.Context, privHex string, nonce uint64, to string, value int64, erc20address string, tokenAmount int) (common.Hash, error) {

	privateKey, err := crypto.HexToECDSA(privHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := pubkeyToAddress(*publicKeyECDSA)
	erc20Instance, err := ct.NewTNT20VoucherContract(common.HexToAddress(erc721address), c.client)
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}

	// gas price
	gasPrice := c.getGasPriceSuggestion(ctx)
	nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	// address
	toAddress := common.HexToAddress(to)

	auth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(privateKey), chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	// auth.Value = big.NewInt(20000000000000000000) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = &gasPrice

	res, err := erc20Instance.TransferFrom(auth, fromAddress, toAddress, big.NewInt(int64(tokenAmount)))

	if err != nil {
		return common.BytesToHash([]byte("")), err
	}

	return res.Hash(), nil
}
func (c EthClient) CrossChainTNT20Transfer(ctx context.Context, privHex string, nonce uint64, to string, value int64, contractAddress string, tokenAmount int) (common.Hash, error) {

	privateKey, err := crypto.HexToECDSA(privHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := pubkeyToAddress(*publicKeyECDSA)
	subchainTNT20Address := common.HexToAddress("0x5C3159dDD2fe0F9862bC7b7D60C1875fa8F81337")
	erc20TokenBank, err := ct.NewTNT20TokenBank(common.HexToAddress(contractAddress), c.client)
	subchainTNT20Instance, _ := ct.NewMockTNT20(subchainTNT20Address, c.client)
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	// gas price
	gasPrice := c.getGasPriceSuggestion(ctx)
	nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	// address
	//toAddress := common.HexToAddress(to)

	auth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(privateKey), chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = common.Big0
	// auth.Value = big.NewInt(20000000000000000000) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = &gasPrice

	subchainTNT20Instance.Approve(auth, common.HexToAddress(contractAddress), big.NewInt(100))

	nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = crossChainFee
	//time.Sleep(1 * time.Second)
	fmt.Println(subchainTNT20Instance.Allowance(nil, fromAddress, common.HexToAddress(contractAddress)))

	fmt.Println()
	res, err := erc20TokenBank.LockTokens(auth, big.NewInt(360888), subchainTNT20Address, fromAddress, big.NewInt(1))
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	fmt.Println(erc20TokenBank.TokenLockNonceMap(nil, big.NewInt(360888)))
	// receipt, err := c.client.TransactionReceipt(context.Background(), res.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if receipt.Status != 1 {
	// 	log.Fatal("lock error")
	// }
	// // fmt.Println("success")
	// fmt.Println(receipt.Logs)
	// fmt.Println(receipt.Logs[2].Data)
	//resolveNum := Resolve(receipt.Logs[2].Data)
	startTime := time.Now()
	mutex.Lock()
	txMapCrossChain[countChainTx1.String()] = startTime
	countChainTx1.Add(big.NewInt(1), countChainTx1)
	mutex.Unlock()
	CountNum += 1
	if CountNum%100 == 0 {
		fmt.Println("already send ", CountNum)
	}
	return res.Hash(), nil
}
func (c EthClient) CrossSubChainTNT20Transfer(ctx context.Context, privHex string, nonce uint64, to string, value int64, contractAddress string, tokenAmount int) (common.Hash, error) {

	privateKey, err := crypto.HexToECDSA(privHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := pubkeyToAddress(*publicKeyECDSA)
	subchainTNT20Address := common.HexToAddress("0x5C3159dDD2fe0F9862bC7b7D60C1875fa8F81337")
	erc20TokenBank, err := ct.NewTNT20TokenBank(common.HexToAddress(contractAddress), c.client)
	subchainTNT20Instance, _ := ct.NewMockTNT20(subchainTNT20Address, c.client)
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	// gas price
	gasPrice := c.getGasPriceSuggestion(ctx)
	nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	// address
	//toAddress := common.HexToAddress(to)

	auth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(privateKey), chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = common.Big0
	// auth.Value = big.NewInt(20000000000000000000) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = &gasPrice

	subchainTNT20Instance.Approve(auth, common.HexToAddress(contractAddress), big.NewInt(100))

	nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = crossChainFee
	//time.Sleep(1 * time.Second)
	fmt.Println(subchainTNT20Instance.Allowance(nil, fromAddress, common.HexToAddress(contractAddress)))

	fmt.Println()
	res, err := erc20TokenBank.LockTokens(auth, big.NewInt(366), subchainTNT20Address, fromAddress, big.NewInt(1))
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	receipt, err := c.client.TransactionReceipt(context.Background(), res.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(receipt.Logs)
	if receipt.Status != 1 {
		log.Fatal("lock error")
	}
	fmt.Println("success")
	startTime := time.Now()
	mutex.Lock()
	txMap[res.Hash()] = startTime
	mutex.Unlock()
	CountNum += 1
	if CountNum%100 == 0 {
		fmt.Println("already send ", CountNum)
	}
	return res.Hash(), nil
}
