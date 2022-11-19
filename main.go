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
	// ThetaRpc = []string{"http://10.10.1.5:16900/rpc", "http://10.10.1.5:16900/rpc", "http://10.10.1.5:16900/rpc", "http://10.10.1.5:16900/rpc"}
	// EthRpc   = []string{"http://10.10.1.5:19888/rpc", "http://10.10.1.5:19888/rpc", "http://10.10.1.5:19888/rpc", "http://10.10.1.5:19888/rpc"} // testnet
	ThetaRpc         = []string{"http://127.0.0.1:16900/rpc", "http://127.0.0.1:16900/rpc", "http://127.0.0.1:16900/rpc", "http://127.0.0.1:16900/rpc"}
	EthRpc           = []string{"http://127.0.0.1:19888/rpc", "http://127.0.0.1:19888/rpc", "http://127.0.0.1:19888/rpc", "http://127.0.0.1:19888/rpc"} // testnet
	Timeout          = 15 * time.Second
	MaxConcurrency   = runtime.NumCPU()
	mesuringDuration = 240 * time.Second //执行数据时间
	queueSize        = 999999            //队列大小
	concurrency      = 1                 //并发数量
	queue            = tps.NewQueue(queueSize)
	closing          uint32
	tpsClosing       uint32
	idlingDuration   uint32
	logLevel         = tps.WARN_LEVEL // INFO_LEVEL, WARN_LEVEL, FATAL_LEVEL
	logger           = tps.NewLogger(logLevel)
	privs            = []string{"57cf79e443d80c5681b5eb44a6e686f8d2289f0b15a784371aa16bbb976780aa",
		"728a2e5396dad4097ada540da0302d4cae3cd75ad94e151a43ec55e1c4c8cc5e",
		"0a709015ad8cd76f66f22b4ede28c5af6a3f4f2e0621d4d41e3e25405a9078a5",
		"1b61f740a6db3648972eb6721b2144f1df5d7851e4eb491ade82640a2b90e704",
		"daab42f8f14cce9b09f1171518d32b62958f92997c9388d06d77def278fbb229",
		"83afc13820ea817f20e543c110505e2d0116aaf27c83297196fdaf954d1465e4",
		"4fa984962825e78281d01074896a111355dd331b97abe358a9ee371afbbf2ccd",
		"142db3c669b0d0e1a34cc5185e1bcd4709cc2175c00d039b537ef85372d634d3",
		"bd07e51c0776035198274aea5589cf81bca664bcde6ed25b068e255a2b9bc8e1",
		"1c9e43a31975347a82b3ddf16fae86b8ead47190b860ff7a5ba93798480ef8a9",
		"f24563e37e24569426bf4d38fa7a0b95e37a2bd9c6336d86310cd9e26a8524e5",
		"c94aef092b9800a35f924bd9d8092717100fc60f6b212b76d44dfcd76491c1d3",
		"41526bca89584202ff2ec68f220781fa63e933b0d5710b1a8c5b94f2b6fbcd7a",
		"c335a892c7cf559780c2888c89906ab5ec19db6e8abc2ba090e4765231b549c6",
		"4e279c0e0b3839398533ebb4aa7b47e9f417bba636c5a780ebcce50dfef9b2a2",
		"d7a0529fc4e96af87cdbccfc74f53c4e7ce42699b524dd5c757db9ff12e196d4"}

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
	client_number    = 1
	clientID         int
	crossPercentage  = 100
)

func hugeBanner(content string, phase int) {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Entering phase ", phase)
	fmt.Println(" ######################################################### ")
	fmt.Println("#                           ", content, "                              #")
	fmt.Println("#  _    _      _ _         _______ _          _           #")
	fmt.Println("#  | |  | |    | | |       |__   __| |        | |         #")
	fmt.Println("#  | |__| | ___| | | ___      | |  | |__   ___| |_ __ _   #")
	fmt.Println("#  |  __  |/ _ \\ | |/ _ \\     | |  | '_ \\ / _ \\ __/ _` |  #")
	fmt.Println("#  | |  | |  __/ | | (_) |    | |  | | | |  __/ || (_| |  #")
	fmt.Println("#  |_|  |_|\\___|_|_|\\___/     |_|  |_| |_|\\___|\\__\\__,_|  #")
	fmt.Println("#                                                         #")
	fmt.Println("#                           ", content, "                              #")
	fmt.Println(" ######################################################### ")
	fmt.Println("")
	fmt.Println("")
}

func main() {
	if len(os.Args) == 3 {
		clientID, _ = strconv.Atoi(os.Args[1])
		model = os.Args[2]
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
		// time.Sleep(mesuringDuration * 2)
		time.Sleep(300 * time.Second)
	}()
	go func() {
		defer hugeBanner("Malicious node is on!", 1)
		// defer fmt.Println("--------------------------------Malicious node is on!-----------------------------------")
		time.Sleep(60 * time.Second)
	}()
	go func() {
		defer hugeBanner("Timing issue is on!", 2)
		// defer fmt.Println("--------------------------------Timing issue is on!--------------------------------")
		time.Sleep(120 * time.Second)
	}()
	go func() {
		defer hugeBanner("Timing issue is off!", 3)
		// defer fmt.Println("--------------------------------Timing issue is off!--------------------------------")
		time.Sleep(180 * time.Second)
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
	if model == "CrossChain" {
		//在跨链测试时需要开一个新的client在另一条链进行监测
		// newclient, err = NewClient("http://127.0.0.1:16888/rpc", "http://127.0.0.1:18888/rpc")
		// newclient, err = NewClient("http://128.110.96.161:16888/rpc", "http://128.110.96.161:18888/rpc")
		newclient, err = NewClient("http://128.110.96.107:16900/rpc", "http://128.110.96.107:19888/rpc") //改真实ip
		// subchain 16900 19888 sidechain "http://127.0.0.1:17900/rpc", "http://127.0.0.1:19988/rpc" mainchain "http://127.0.0.1:16888/rpc", "http://127.0.0.1:18888/rpc"
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
