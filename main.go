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
)

const (
	ERC721 = "erc721"
	EHT    = "eth"
	Cross  = "CrossChainTNT20"
)

var (
	ThetaRpc         = []string{"http://127.0.0.1:16888/rpc", "http://20.163.221.185:16888/rpc", "http://20.220.200.72:16888/rpc", "http://20.231.77.191:16888/rpc"}
	EthRpc           = []string{"http://127.0.0.1:18888/rpc", "http://127.0.0.1:18888/rpc"} // testnet
	Timeout          = 15 * time.Second
	MaxConcurrency   = runtime.NumCPU()
	mesuringDuration = 120 * time.Second //执行数据时间
	queueSize        = 999999            //队列大小
	concurrency      = 1                 //并发数量
	queue            = tps.NewQueue(queueSize)
	closing          uint32
	tpsClosing       uint32
	idlingDuration   uint32
	logLevel         = tps.WARN_LEVEL // INFO_LEVEL, WARN_LEVEL, FATAL_LEVEL
	logger           = tps.NewLogger(logLevel)
	privs            = []string{
		"a249a82c42a282e87b2ddef63404d9dfcf6ea501dcaf5d447761765bd74f666d",
		"93a90ea508331dfdf27fb79757d4250b4e84954927ba0073cd67454ac432c737",
		"d0d53ac0b4cd47d0ce0060dddc179d04145fea2ee2e0b66c3ee1699c6b492013",
		"83f0bb8655139cef4657f90db64a7bb57847038a9bd0ccd87c9b0828e9cbf76d",
		"8888888888888888888888888888888888888888888888888888888888888888",
	}

	model = "Theta" //压测类型

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
	countChainTx1    = big.NewInt(0)
	countChainTx2    = big.NewInt(0)
	txMapCrossChain  map[string]time.Time
	client_number    = concurrency
)

func main() {
	txMap = make(map[common.Hash]time.Time)
	txMapCrossChain = make(map[string]time.Time)
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
	var client_list []EthClient
	var err error
	for i := 0; i < client_number; i++ {
		var client EthClient

		if model == "CrossChainTNT20" {
			client, err = NewClient("http://localhost:17900/rpc",EthRpc[i])
		} else {
			client, err = NewClient(ThetaRpc[i],"")
		}

		if err != nil {
			logger.Fatal("err NewClient: ", err)
		}
		client_list = append(client_list, client)
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
		ethStressTest(&client_list, ctx)
	}
	var newclient EthClient
	if model == "CrossChainTNT20" {
		newclient, err = NewClient("http://localhost:17900/rpc",EthRpc[0])
		//newclient.client, err = ethclient.Dial("http://localhost:19988/rpc")
	} else {
		newclient = client_list[0]
	}
	fmt.Println("-----------Start Measuring----------")
	if err = tps.StartTPSMeasuring(context.Background(), &newclient, &tpsClosing, &idlingDuration, logger); err != nil {
		fmt.Println("err StartTPSMeasuring:", err)
		logger.Fatal("err StartTPSMeasuring: ", err)
	}

}
