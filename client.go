package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
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
				if model == "CrossChainTNT20" {
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
						mutex.Lock()

						startTime, ok := txMapCrossChain[event.VoucherMintNonce.String()]
						if ok {
							countChainTx2.Add(event.VoucherMintNonce, big.NewInt(1))
							mutex.Unlock()
							elapsedTime = elapsedTime + time.Since(startTime)/time.Millisecond
							fmt.Println("latency is", time.Since(startTime))
							result += 1
						} else {
							fmt.Println("unfind!", event.VoucherMintNonce)
						}
					}

				} else if model == "CrossSubChainTNT20" {
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
	txMap[common.HexToHash(result.TxHash).String()] = startTime
	mutex.Unlock()
	CountNum += 1
	// if CountNum%100 == 0 {
	// 	fmt.Println("already send ", CountNum)
	// }
	fmt.Println("have send ", CountNum, " txs", nonce+1)
	return common.Hash{}, err //common.BytesToHash(formatted), err
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
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = common.Big0
	// auth.Value = big.NewInt(20000000000000000000) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = &gasPrice

	_, err = subchainTNT20Instance.Approve(auth, common.HexToAddress(contractAddress), big.NewInt(100))
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	//nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
	auth.Nonce = big.NewInt(int64(nonce + 1))
	auth.Value = crossChainFee
	//time.Sleep(1 * time.Second)
	//fmt.Println(subchainTNT20Instance.Allowance(nil, fromAddress, common.HexToAddress(contractAddress)))

	//fmt.Println()
	// res, err := subchainTNT20Instance.Mint(auth, fromAddress, big.NewInt(1000000000))
	// if err != nil {
	// 	return common.BytesToHash([]byte("")), err
	// }
	res, err := erc20TokenBank.LockTokens(auth, big.NewInt(360777), subchainTNT20Address, fromAddress, big.NewInt(1))
	if err != nil {
		return common.BytesToHash([]byte("")), err
	}
	lockNonce, _ := erc20TokenBank.TokenLockNonceMap(nil, big.NewInt(360777))
	fmt.Println("lock", lockNonce)
	startTime := time.Now()
	mutex.Lock()
	txMapCrossChain[lockNonce.String()] = startTime
	mutex.Unlock()
	receipt, err := c.client.TransactionReceipt(context.Background(), res.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		log.Fatal("lock error")
	}
	fmt.Println("success")
	// fmt.Println(receipt.Logs)
	// fmt.Println(receipt.Logs[2].Data)
	//resolveNum := Resolve(receipt.Logs[2].Data)

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
	time.Sleep(50 * time.Millisecond)
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
