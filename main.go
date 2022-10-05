package main

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blockchain-tps-test/samples/theta/tps"
	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/thetasubchain/eth/ethclient"
)

const (
	ERC721 = "erc721"
	EHT    = "eth"
	Cross  = "CrossChainTNT20"
)

var (
	Endpoint         = "http://localhost:16888/rpc" // testnet
	Timeout          = 15 * time.Second
	MaxConcurrency   = runtime.NumCPU()
	mesuringDuration = 120 * time.Second //执行数据时间
	queueSize        = 10000000          //队列大小
	concurrency      = 1                 //并发数量
	queue            = tps.NewQueue(queueSize)
	closing          uint32
	tpsClosing       uint32
	idlingDuration   uint32
	logLevel         = tps.WARN_LEVEL // INFO_LEVEL, WARN_LEVEL, FATAL_LEVEL
	logger           = tps.NewLogger(logLevel)
	privs            = []string{
		"93a90ea508331dfdf27fb79757d4250b4e84954927ba0073cd67454ac432c737",
		"1111111111111111111111111111111111111111111111111111111111111111",
		"4444444444444444444444444444444444444444444444444444444444444444",
		"7777777777777777777777777777777777777777777777777777777777777777",
		"8888888888888888888888888888888888888888888888888888888888888888",
	}

	model = Cross //压测类型

	addr_priv     = make(map[string]string, len(privs))
	erc721address = "0x0000000000000000000000000000000000000009"
	client        EthClient
	txMap         map[common.Hash]time.Time

	avgLatency       time.Duration
	mutex            sync.Mutex
	CountNum         int
	chainID          = big.NewInt(360777)
	Erc20Address     = ""
	TokenBankAddress = "0x47e9Fbef8C83A1714F1951F142132E6e90F5fa5D"
	countChainTx1    = 0
	countChainTx2    = 0
	txMapCrossChain  map[int]time.Time
)

func main() {
	txMap = make(map[common.Hash]time.Time)
	txMapCrossChain = make(map[int]time.Time)
	// client := jsonrpc.NewRPCClient("http://localhost:16888/rpc")
	// // res, err := client1.BlockNumber(context.Background())
	// rpcRes, rpcErr := client.Call("theta.GetBlockByHeight", trpc.GetBlockByHeightArgs{
	// 	Height: tcommon.JSONUint64(378)})

	// //logger.Infof("eth_getBlockByNumber, rpcRes: %v, rpcRes.Rsult: %v", rpcRes, rpcRes.Result)
	// chainID := new(big.Int)
	// chainID.SetString("privatenet", 16)
	// result, err := ethrpc.GetBlockFromTRPCResult(chainID, rpcRes, rpcErr, true)

	// fmt.Println(err, result.Transactions)
	// client, err := ethclient.Dial("http://localhost:18888/rpc")
	// fmt.Println(err)
	// block, err:= client.BlockByNumber(context.Background(), big.NewInt(int64(378)))
	// fmt.Println(err)
	// fmt.Println(len(block.Transactions()))
	go func() {
		//停止发送交易时间
		defer atomic.AddUint32(&closing, 1)
		time.Sleep(mesuringDuration)
	}()

	go func() {
		//统计tps结束时间
		defer atomic.AddUint32(&tpsClosing, 1)
		time.Sleep(mesuringDuration * 2)
	}()

	var client EthClient
	var err error
	if model == "CrossChainTNT20" {
		client, err = NewClient("http://localhost:17900/rpc")
	} else {
		client, err = NewClient(Endpoint)
	}

	if err != nil {
		logger.Fatal("err NewClient: ", err)
	}
	//fmt.Println(client.CountTx(context.Background(), 3847))
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	if model == "ERC20" {
		erc20StressTest(&client, ctx)
	} else if model == "CrossChainTNT20" {
		crossChainTNT20StressTest(&client, ctx)
	} else if model == "CrossSubChainTNT20" {
		crossSubChainTNT20StressTest(&client, ctx)
	} else {
		ethStressTest(&client, ctx)
	}
	var newclient EthClient
	if model == "CrossChainTNT20" {
		newclient, err = NewClient("http://localhost:17900/rpc")
		newclient.client, err = ethclient.Dial("http://localhost:19988/rpc")
	}
	fmt.Println("-----------Start Measuring----------")
	if err = tps.StartTPSMeasuring(context.Background(), &newclient, &tpsClosing, &idlingDuration, logger); err != nil {
		fmt.Println("err StartTPSMeasuring:", err)
		logger.Fatal("err StartTPSMeasuring: ", err)
	}

}
