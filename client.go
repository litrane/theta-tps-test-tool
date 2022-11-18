package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/blockchain-tps-test/samples/theta/tps"
	"github.com/pkg/errors"
	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/common/hexutil"

	"github.com/thetatoken/theta/crypto"
	"github.com/thetatoken/theta/ledger/types"
	"github.com/thetatoken/theta/rpc"
	"github.com/ybbus/jsonrpc"

	ct "github.com/blockchain-tps-test/samples/theta/accessors"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/thetasubchain/eth/abi"
	"github.com/thetatoken/thetasubchain/eth/abi/bind"
	"github.com/thetatoken/thetasubchain/eth/ethclient"
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

func NewClient(rpcClientUrl, ethClientUrl string) (c EthClient, err error) { //,theta-url,eth-url
	c.client, err = ethclient.Dial(ethClientUrl)
	if err != nil {
		return
	}
	c.rpcClient = jsonrpc.NewRPCClient(rpcClientUrl)
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
				if model == "CrossChain" {
					var test RPCResult
					fmt.Println(string(value["receipt"]))
					json.Unmarshal(value["receipt"], &test)
					if len(test.Result) != 0 {
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
						fmt.Println("find", event.VoucherMintNonce)
						fmt.Println("find", event.Denom)
						mutex.Lock()

						startTime, ok := txMapCrossChain[event.VoucherMintNonce.String()]
						if ok {
							countChainTx2.Add(event.VoucherMintNonce, big.NewInt(1))
							mutex.Unlock()
							elapsedTime = elapsedTime + time.Since(startTime)/time.Millisecond
							fmt.Println("latency is", time.Since(startTime))
							fmt.Println("cur time is ", time.Now().UnixNano())
							result += 1
						} else {
							fmt.Println("unfind!", event.VoucherMintNonce)
						}
					}

				} else if model == "Inchain" {
					mutex.Lock()
					//fmt.Println(common.value["hash"].String())
					//var hash string
					//json.Unmarshal(value["raw"], &hash)
					fmt.Println(string(value["raw"]))
					pattern := regexp.MustCompile(`"sequence": "(\d+)"`)
					numberStrings := pattern.FindAllStringSubmatch(string(value["raw"]), -1)
					numbers := make([]int, len(numberStrings))
					for i, numberString := range numberStrings {
						number, err := strconv.Atoi(numberString[1])
						if err != nil {
							panic(err)
						}
						numbers[i] = number
					}

					pattern = regexp.MustCompile(`"address": "(.*?)"`)
					numberStrings = pattern.FindAllStringSubmatch(string(value["raw"]), -1)
					strings1 := make([]string, len(numberStrings))
					for i, numberString := range numberStrings {
						strings1[i] = numberString[1]
						fmt.Println(numberString)
					}
					startTime, ok := txMap[strings.ToUpper(strings1[0])+fmt.Sprint(numbers[0])]
					if ok {
						//fmt.Println(trpcResult.Txs[i].Hash)
						fmt.Println("start ", startTime, "end-start", time.Since(startTime))
						mutex.Unlock()
						elapsedTime = elapsedTime + (time.Since(startTime) / time.Millisecond)
						//fmt.Println(time.Since(startTime) / time.Second)
						//fmt.Println(elapsedTime)
						result += 1
					} else {
						fmt.Println("unfind!", numbers[0])
					}

				}

			} else if types.TxType(trpcResult.Txs[i].Type) == types.TxSend {
				mutex.Lock()
				startTime := txMap[trpcResult.Txs[i].Hash.String()]
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
	//startTime := time.Now()
	rpcResult, err := c.rpcClient.Call("theta.GetBlockByHeight", rpc.GetBlockByHeightArgs{
		Height: common.JSONUint64(height)})
	if err != nil {
		return 0, time.Duration(0), err
	}
	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(rpcResult.Result, "", "    ")

	//logger.Infof("HandleThetaRPCResponse, jsonBytes: %v", strin(jsonBytes))
	result, avg_latency, err := parse(jsonBytes)
	//totalTime := time.Since(startTime) / time.Millisecond
	//fmt.Println("call and parse consume ", totalTime)
	return result, avg_latency, nil
}

func (c EthClient) CountPendingTx(ctx context.Context) (int, error) {

	// count, err := c.client.PendingTransactionCount(ctx)
	// if err != nil {
	// 	return 0, err
	// }
	return int(0), nil
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
	return *gasPrice

}
func (c *EthClient) GetTransactionReceipt(ctx context.Context, hashStr string) (EthGetReceiptResult, error) {
	//logger.Infof("eth_getTransactionReceipt called, txHash: %v", hashStr)

	result := EthGetReceiptResult{}

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := rpc.GetTransactionResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		var objmap map[string]json.RawMessage
		json.Unmarshal(jsonBytes, &objmap)
		if objmap["transaction"] != nil {
			if types.TxType(trpcResult.Type) == types.TxSend {
				tx := types.SendTx{}
				json.Unmarshal(objmap["transaction"], &tx)
				result.From = tx.Inputs[0].Address
				result.To = tx.Outputs[0].Address
			}
			// if types.TxType(trpcResult.Type) == types.TxSmartContract {
			// 	tx := types.SmartContractTx{}
			// 	json.Unmarshal(objmap["transaction"], &tx)
			// 	result.From = tx.From.Address
			// 	result.To = tx.To.Address
			// 	result.ContractAddress = trpcResult.Receipt.ContractAddress
			// }
		}
		return trpcResult, nil
	}

	var thetaGetTransactionResult rpc.GetTransactionResult
	//maxRetry := 5
	for i := 0; i >= 0; i++ { // It might take some time for a tx to be finalized, retry a few times
		rpcRes, rpcErr := c.rpcClient.Call("theta.GetTransaction", rpc.GetTransactionArgs{Hash: hashStr})

		resultIntf, err := HandleThetaRPCResponse(rpcRes, rpcErr, parse)
		if err != nil {
			resultMsg := ""
			if resultIntf != nil {
				resultMsg = resultIntf.(string)
			}
			fmt.Printf("eth_getTransactionReceipt, err: %v, result: %v", err, resultMsg)
			return result, err
		}

		thetaGetTransactionResult = resultIntf.(rpc.GetTransactionResult)
		if thetaGetTransactionResult.Status == rpc.TxStatusFinalized {
			break
		}
		//time.Sleep(1 * time.Second) // one block duration
	}
	if thetaGetTransactionResult.Receipt == nil {
		return result, nil
	}

	result.BlockHash = thetaGetTransactionResult.BlockHash
	result.BlockHeight = hexutil.Uint64(thetaGetTransactionResult.BlockHeight)
	result.TxHash = thetaGetTransactionResult.TxHash
	result.GasUsed = hexutil.Uint64(thetaGetTransactionResult.Receipt.GasUsed)
	result.Logs = make([]EthLogObj, len(thetaGetTransactionResult.Receipt.Logs))
	for i, log := range thetaGetTransactionResult.Receipt.Logs {
		result.Logs[i] = ThetaLogToEthLog(log)
		result.Logs[i].BlockHash = result.BlockHash
		result.Logs[i].BlockHeight = result.BlockHeight
		result.Logs[i].TxHash = result.TxHash
		result.Logs[i].LogIndex = hexutil.Uint64(i)
	}

	//TODO: handle logIndex & TransactionIndex of logs

	if thetaGetTransactionResult.Receipt.EvmErr == "" {
		result.Status = 1
	} else {
		result.Status = 0
	}

	return result, nil
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
	subchainTNT20Address := common.HexToAddress("0x47c5e40890bcE4a473A49D7501808b9633F29782") //subchain 0x5C3159dDD2fe0F9862bC7b7D60C1875fa8F81337 mainchain 0x59AF421cB35fc23aB6C8ee42743e6176040031f4
	erc20TokenBank, err := ct.NewTNT20TokenBank(common.HexToAddress(contractAddress), c.client)
	subchainTNT20Instance, _ := ct.NewMockTNT20(subchainTNT20Address, c.client)
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	// gas price
	gasPrice := c.getGasPriceSuggestion(ctx)
	//nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	// address
	//toAddress := common.HexToAddress(to)
	rr, rt := subchainTNT20Instance.BalanceOf(nil, fromAddress)
	fmt.Println("coin is ", rr, rt)
	auth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(privateKey), chainID)
	if err != nil {
		log.Fatal(err)
	}
	// auth.Value = big.NewInt(20000000000000000000) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = &gasPrice

	//nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = crossChainFee
	//time.Sleep(1 * time.Second)
	//fmt.Println(subchainTNT20Instance.Allowance(nil, fromAddress, common.HexToAddress(contractAddress)))

	//fmt.Println()
	// res, err := subchainTNT20Instance.Mint(auth, fromAddress, big.NewInt(1000000000))
	// if err != nil {
	// 	return common.BytesToHash([]byte("")), err
	// }
	res, err := erc20TokenBank.LockTokens(auth, big.NewInt(366), subchainTNT20Address, fromAddress, big.NewInt(1))
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	// time.Sleep(1 * time.Second)
	receipt, err := c.GetTransactionReceipt(context.Background(), res.Hash().Hex())
	if err != nil {
		log.Fatal(err)
	}

	if receipt.Status != 1 {
		log.Fatal("lock error")
	}
	lockNonce, _ := erc20TokenBank.TokenLockNonceMap(nil, big.NewInt(366))
	fmt.Println("lock", lockNonce)
	startTime := time.Now()
	mutex.Lock()
	txMapCrossChain[lockNonce.String()] = startTime
	mutex.Unlock()
	fmt.Println("success")
	// fmt.Println(receipt.Logs)
	// fmt.Println(receipt.Logs[2].Data)
	//resolveNum := Resolve(receipt.Logs[2].Data)
	time.Sleep(1 * time.Second)
	CountNum += 1
	if CountNum%100 == 0 {
		fmt.Println("already send ", CountNum)
	}
	return res.Hash(), nil
}

func (c EthClient) CrossSubChainTNT20Transfer(ctx context.Context, privHex string, nonce uint64, to string, value int64, contractAddress string, tokenAmount int) (common.Hash, error) {

	//fmt.Println("send1", nonce+1, "send2", nonce+2)
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
	subchainTNT20Address := common.HexToAddress("0x5C3159dDD2fe0F9862bC7b7D60C1875fa8F81337") // subchain 0x5C3159dDD2fe0F9862bC7b7D60C1875fa8F81337 mainchain 0x59AF421cB35fc23aB6C8ee42743e6176040031f4
	erc20TokenBank, err := ct.NewTNT20TokenBank(common.HexToAddress(contractAddress), c.client)
	//subchainTNT20Instance, _ := ct.NewMockTNT20(subchainTNT20Address, c.client)
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	// gas price
	gasPrice := c.getGasPriceSuggestion(ctx)
	//nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
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

	//_, err1 := subchainTNT20Instance.Approve(auth, common.HexToAddress(contractAddress), big.NewInt(99999999))
	//fmt.Println(err1)
	//time.Sleep(50 * time.Millisecond)
	//nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	//auth.Nonce = big.NewInt(int64(nonce + 1))
	auth.Value = crossChainFee
	//time.Sleep(1 * time.Second)
	//fmt.Println(subchainTNT20Instance.Allowance(nil, fromAddress, common.HexToAddress(contractAddress)))

	//fmt.Println()
	res, err := erc20TokenBank.LockTokens(auth, big.NewInt(360777), subchainTNT20Address, fromAddress, big.NewInt(1))
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	//time.Sleep(50 * time.Millisecond)
	// receipt, err := c.client.TransactionReceipt(context.Background(), res.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if receipt.Status != 1 {
	// 	log.Fatal("lock error")
	// }
	// fmt.Println("success")
	startTime := time.Now()
	//mutex.Lock()
	//fmt.Println(res.Hash().String())
	txMap[strings.ToUpper(fromAddress.Hex())+fmt.Sprint(nonce+1)] = startTime
	//mutex.Unlock()
	CountNum += 1
	if CountNum%100 == 0 {
		fmt.Println("already send ", CountNum)
	}
	fmt.Println(fromAddress.Hex())
	return res.Hash(), nil
}
