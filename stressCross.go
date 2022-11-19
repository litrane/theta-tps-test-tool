package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/blockchain-tps-test/samples/theta/tps"
	"github.com/pkg/errors"
	"github.com/thetatoken/theta/crypto"
)

// func crossChainTNT20StressTest(client *[]EthClient, ctx context.Context) {
// //初始化钱包，首先根据私钥生成公钥地址addrs
// 	addrs := make([]string, len(privs))
// 	for i := range privs {
// 		privateKey, err := crypto.HexToECDSA(privs[i])

// 		// privateKey, err := crypto.HexToECDSA("2dad160420b1e9b6fc152cd691a686a7080a0cee41b98754597a2ce57cc5dab1")
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		publicKey := privateKey.Public()
// 		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 		if !ok {
// 			log.Fatal("error casting public key to ECDSA")
// 		}

// 		fromAddress := pubkeyToAddress(*publicKeyECDSA)
// 		addrs[i] = fromAddress.Hex()
// 	}
// //生成钱包，第一次new时就获取account的nonce，之后每次发送交易时nonce+1
// 	var wallet_list []tps.Wallet
// 	for i := 0; i < client_number; i++ {
// 		wallet_single, err := tps.NewWallet(ctx, (*client)[i], privs, addrs)
// 		if err != nil {
// 			logger.Fatal("err NewWallet: ", err)
// 		}
// 		wallet_list = append(wallet_list, wallet_single)
// 	}
// 	var err error
// //生成worker需要不停循环的task
// 	taskDo := func(t tps.Task, id int) error {
// 		task, ok := t.(*EthTask)
// 		if !ok {
// 			return errors.New("unexpected task type")
// 		}

// 		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
// 		defer cancel()
// //根据id选择对应的钱包以及密钥
// 		var (
// 			priv         = wallet_list[id].Priv(id)
// 			currentNonce = wallet_list[id].CurrentNonce(priv)
// 		)
// 		err = task.Do(ctx, &(*client)[id], priv, currentNonce, &queue, logger, TokenBankAddress)
// 		//因为一次do里面包含了两笔交易，一个是approve（这个日后可以移至init），另一个是lock，所以需要增加两次nonce
// 		wallet_list[id].IncrementNonce(priv)
// 		wallet_list[id].IncrementNonce(priv)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "ValidateInputAdvanced: Got") {//错误捕捉，根据返回期盼的nonce和实际的nonce进行比较，如果不一致则重新设置nonce
// 				pattern := regexp.MustCompile(`(\d+)`)
// 				numberStrings := pattern.FindAllStringSubmatch(err.Error(), -1)
// 				numbers := make([]int, len(numberStrings))
// 				for i, numberString := range numberStrings {
// 					number, err := strconv.Atoi(numberString[1])
// 					if err != nil {
// 						panic(err)
// 					}
// 					numbers[i] = number
// 				}
// 				wallet_list[id].RecetNonce(priv, uint64(numbers[3]))
// 				fmt.Println("Restnonce is", wallet_list[id].CurrentNonce(priv))
// 				return nil
// 			}
// 			if errors.Is(err, tps.ErrWrongNonce) {
// 				wallet_list[id].RecetNonce(priv, currentNonce)
// 				task.tryCount = 0
// 				queue.Push(task)
// 				return nil
// 			}
// 			if errors.Is(err, tps.ErrTaskRetry) {
// 				wallet_list[id].IncrementNonce(priv)
// 				return nil
// 			}
// 			return errors.Wrap(err, "err Do")
// 		}

// 		return nil
// 	}
// //将需要做的task给worker实例化
// 	worker := tps.NewWorker(taskDo)
// //根据并发数启动worker，每个worker用的client链接都不一样
// 	if concurrency > MaxConcurrency {
// 		logger.Warn(fmt.Sprintf("concurrency setting is over logical max(%d)", MaxConcurrency))
// 	}
// 	for i := 0; i < concurrency; i++ {
// 		go worker.Run(&queue, i)
// 	}
// //不停地往队列中添加task，但是这个逻辑目前没有太用到，因为task的信息没有用到，但是可以指定一共发送多少笔交易停止
// 	go func() {
// 		count := 2
// 		for {
// 			if atomic.LoadUint32(&closing) == 1 {
// 				break
// 			}

// 			if queue.CountTasks() > queueSize {
// 				continue
// 			}

// 			queue.Push(&EthTask{
// 				to:      "0x27F6F1bb3e2977c3CB014e7d4B5639bB133A6032",
// 				amount:  1,
// 				tokenId: int64(count),
// 			})
// 			count++
// 		}
// 	}()
// }
//下面的函数同理，因为要发送不同类型的交易所以需要不同的stress函数，但是可能可以合并成一个函数？这个工作还没有做，因为需要统一错误捕捉
func crossSubChainTNT20StressTest(client *[]EthClient, ctx context.Context) {
	privsGroup := [][]string{}
	addrs := [][]string{}
	for key, value := range privs {
		// if(key==0){
		// 	continue;
		// }
		if len(privsGroup) < key/(len(privs)/client_number)+1 {
			privsGroup = append(privsGroup, []string{})
		}
		if len(addrs) < key/(len(privs)/client_number)+1 {
			addrs = append(addrs, []string{})
		}
		privsGroup[key/(len(privs)/client_number)] = append(privsGroup[key/(len(privs)/client_number)], value)
		privateKey, err := crypto.HexToECDSA(value)

		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		fromAddress := pubkeyToAddress(*publicKeyECDSA)
		addrs[key/(len(privs)/client_number)] = append(addrs[key/(len(privs)/client_number)], fromAddress.Hex())
	}

	var wallet_list []tps.Wallet
	for i := 0; i < client_number; i++ {
		wallet_single, err := tps.NewWallet(ctx, (*client)[i], privsGroup[i], addrs[i])
		if err != nil {
			logger.Fatal("err NewWallet: ", err)
		}
		wallet_list = append(wallet_list, wallet_single)
	}

	//var err error
	taskDo := func(t tps.Task, id int) error {
		task, ok := t.(*EthTask)
		if !ok {
			return errors.New("unexpected task type")
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var (
			priv         = wallet_list[id].Priv(int(wallet_list[id].SendCount()) % wallet_list[id].PrivsLength())
			currentNonce = wallet_list[id].CurrentNonce(priv)
		)
		_ = task.Do(ctx, &(*client)[id], priv, currentNonce, &queue, logger, TokenBankAddress)
		wallet_list[id].IncrementNonce(priv)
		//wallet_list[id].IncrementNonce(priv)
		wallet_list[id].IncrementSendCount()
		// if err != nil {
		// 	if strings.Contains(err.Error(), "ValidateInputAdvanced: Got") {
		// 		pattern := regexp.MustCompile(`(\d+)`)
		// 		numberStrings := pattern.FindAllStringSubmatch(err.Error(), -1)
		// 		numbers := make([]int, len(numberStrings))
		// 		for i, numberString := range numberStrings {
		// 			number, err := strconv.Atoi(numberString[1])
		// 			if err != nil {
		// 				panic(err)
		// 			}
		// 			numbers[i] = number
		// 		}
		// 		wallet_list[id].RecetNonce(priv, uint64(numbers[3]))
		// 		fmt.Println("Restnonce is", wallet_list[id].CurrentNonce(priv))
		// 		return nil
		// 	}
		// 	if errors.Is(err, tps.ErrWrongNonce) {
		// 		wallet_list[id].RecetNonce(priv, currentNonce)
		// 		fmt.Println("Restnonce is", currentNonce)
		// 		task.tryCount = 0
		// 		queue.Push(task)
		// 		return nil
		// 	}
		// 	if errors.Is(err, tps.ErrTaskRetry) {
		// 		wallet_list[id].IncrementNonce(priv)
		// 		fmt.Println("IncrementNonce is", currentNonce)
		// 		return nil
		// 	}
		// 	return errors.Wrap(err, "err Do")
		// }

		return nil
	}

	worker := tps.NewWorker(taskDo)

	if concurrency > MaxConcurrency {
		logger.Warn(fmt.Sprintf("concurrency setting is over logical max(%d)", MaxConcurrency))
	}

	go worker.Run(&queue, clientID)

	go func() {
		count := 2
		for {
			if atomic.LoadUint32(&closing) == 1 {
				break
			}

			if queue.CountTasks() > queueSize {
				continue
			}
			if count == 10000 {
				fmt.Println("have send 1000")
				break
			}
			// if count%crossPercentage == 0 {
			// 	queue.Push(&EthTask{
			// 		to:            "0x27F6F1bb3e2977c3CB014e7d4B5639bB133A6032",
			// 		amount:        1,
			// 		tokenId:       int64(count),
			// 		transfer_type: "CrossChain",
			// 	})
			// } else {
			queue.Push(&EthTask{
				to:            "0x27F6F1bb3e2977c3CB014e7d4B5639bB133A6032",
				amount:        1,
				tokenId:       int64(count),
				transfer_type: "InChain",
			})
			// }
			count++
		}
	}()
}
