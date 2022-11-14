package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blockchain-tps-test/samples/theta/tps"
)

const (
	ERC721 = "erc721"
	EHT    = "eth"
	Cross  = "CrossChainTNT20"
)

var (
	ThetaRpc         = []string{"http://10.10.1.5:16900/rpc", "http://127.0.0.1:16888/rpc", "http://20.220.200.72:16888/rpc", "http://20.231.77.191:16888/rpc"}
	EthRpc           = []string{"http://10.10.1.5:19888/rpc", "http://127.0.0.1:18888/rpc"} // testnet
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
		"472d4146a9fb59433d76c42be9ff0d8a9e1cfbaaff01cbdf68e858715c85af01",
		"a249a82c42a282e87b2ddef63404d9dfcf6ea501dcaf5d447761765bd74f666d",
		"d0d53ac0b4cd47d0ce0060dddc179d04145fea2ee2e0b66c3ee1699c6b492013",
		"83f0bb8655139cef4657f90db64a7bb57847038a9bd0ccd87c9b0828e9cbf76d",
		"8888888888888888888888888888888888888888888888888888888888888888",
	}

	model = "CrossChainTNT20" //压测类型

	addr_priv     = make(map[string]string, len(privs))
	erc721address = "0x0000000000000000000000000000000000000009"
	client        EthClient
	txMap         map[string]time.Time

	avgLatency       time.Duration
	mutex            sync.Mutex
	CountNum         int
	chainID          = big.NewInt(360777) // 366 360777
	Erc20Address     = ""
	TokenBankAddress = "0x47e9Fbef8C83A1714F1951F142132E6e90F5fa5D" // subchain 0x47e9Fbef8C83A1714F1951F142132E6e90F5fa5D mainchain 0x2Ce636d6240f8955d085a896e12429f8B3c7db26
	countChainTx1    = big.NewInt(0)
	countChainTx2    = big.NewInt(0)
	txMapCrossChain  map[string]time.Time
	client_number    = concurrency
	clientID         int
	crossPercentage  = 100
)

func main() {
	if len(os.Args) == 1 {
		clientID, _ = strconv.Atoi(os.Args[0])
	} else {
		fmt.Println("Wrong Input Arguments!")
	}
	txMap = make(map[string]time.Time)
	txMapCrossChain = make(map[string]time.Time)
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
	//初始化客户端
	for i := 0; i < client_number; i++ {
		var client EthClient

		if model == "CrossSubChainTNT20" {
			client, err = NewClient(ThetaRpc[i], EthRpc[i])
		} else {
			client, err = NewClient(ThetaRpc[i], EthRpc[i])
		}

		if err != nil {
			logger.Fatal("err NewClient: ", err)
		}
		client_list = append(client_list, client)
	}
	//还未测试完init token
	//init_token(client_list[0], privs)
	//return
	//开始进行压测
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	crossSubChainTNT20StressTest(&client_list, ctx)

	var newclient EthClient
	if model == "CrossChainTNT20" {
		//在跨链测试时需要开一个新的client在另一条链进行监测
		newclient, err = NewClient("http://10.10.1.1:16888/rpc", "http://10.10.1.1:18888/rpc") // subchain 16900 19888 sidechain "http://127.0.0.1:17900/rpc", "http://127.0.0.1:19988/rpc" mainchain "http://127.0.0.1:16888/rpc", "http://127.0.0.1:18888/rpc"
	} else {
		//否则就用第一个client监测
		newclient = client_list[0]
	}
	//开始TPS以及延迟测量
	fmt.Println("-----------Start Measuring----------")
	if err = tps.StartTPSMeasuring(context.Background(), &newclient, &tpsClosing, &idlingDuration, logger); err != nil {
		fmt.Println("err StartTPSMeasuring:", err)
		logger.Fatal("err StartTPSMeasuring: ", err)
	}

}
