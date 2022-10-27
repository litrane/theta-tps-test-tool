package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/blockchain-tps-test/samples/theta/tps"
	"github.com/pkg/errors"
	"github.com/thetatoken/theta/crypto"
)

func crossChainTNT20StressTest(client *[]EthClient, ctx context.Context) {

	addrs := make([]string, len(privs))
	for i := range privs {
		privateKey, err := crypto.HexToECDSA(privs[i])

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
		addrs[i] = fromAddress.Hex()
	}

	var wallet_list []tps.Wallet
	for i := 0; i < client_number; i++ {
		wallet_single, err := tps.NewWallet(ctx, (*client)[i], privs, addrs)
		if err != nil {
			logger.Fatal("err NewWallet: ", err)
		}
		wallet_list = append(wallet_list, wallet_single)
	}
	var err error

	taskDo := func(t tps.Task, id int) error {
		task, ok := t.(*EthTask)
		if !ok {
			return errors.New("unexpected task type")
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var (
			priv         = wallet_list[id].Priv(id)
			currentNonce = wallet_list[id].CurrentNonce(priv)
		)
		err = task.Do(ctx, &(*client)[id], priv, currentNonce, &queue, logger, TokenBankAddress)
		wallet_list[id].IncrementNonce(priv)
		wallet_list[id].IncrementNonce(priv)
		if err != nil {
			if strings.Contains(err.Error(), "ValidateInputAdvanced: Got") {
				pattern := regexp.MustCompile(`(\d+)`)
				numberStrings := pattern.FindAllStringSubmatch(err.Error(), -1)
				numbers := make([]int, len(numberStrings))
				for i, numberString := range numberStrings {
					number, err := strconv.Atoi(numberString[1])
					if err != nil {
						panic(err)
					}
					numbers[i] = number
				}
				wallet_list[id].RecetNonce(priv, uint64(numbers[3]))
				fmt.Println("Restnonce is", wallet_list[id].CurrentNonce(priv))
				return nil
			}
			if errors.Is(err, tps.ErrWrongNonce) {
				wallet_list[id].RecetNonce(priv, currentNonce)
				task.tryCount = 0
				queue.Push(task)
				return nil
			}
			if errors.Is(err, tps.ErrTaskRetry) {
				wallet_list[id].IncrementNonce(priv)
				return nil
			}
			return errors.Wrap(err, "err Do")
		}

		return nil
	}

	worker := tps.NewWorker(taskDo)

	if concurrency > MaxConcurrency {
		logger.Warn(fmt.Sprintf("concurrency setting is over logical max(%d)", MaxConcurrency))
	}
	for i := 0; i < concurrency; i++ {
		go worker.Run(&queue, i)
	}

	go func() {
		count := 2
		for {
			if atomic.LoadUint32(&closing) == 1 {
				break
			}

			if queue.CountTasks() > queueSize {
				continue
			}

			queue.Push(&EthTask{
				to:      "0x27F6F1bb3e2977c3CB014e7d4B5639bB133A6032",
				amount:  1,
				tokenId: int64(count),
			})
			count++
		}
	}()
}
func crossSubChainTNT20StressTest(client *[]EthClient, ctx context.Context) {

	addrs := make([]string, len(privs))
	for i := range privs {
		privateKey, err := crypto.HexToECDSA(privs[i])

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
		addrs[i] = fromAddress.Hex()
	}

	var wallet_list []tps.Wallet
	for i := 0; i < client_number; i++ {
		wallet_single, err := tps.NewWallet(ctx, (*client)[i], privs, addrs)
		if err != nil {
			logger.Fatal("err NewWallet: ", err)
		}
		wallet_list = append(wallet_list, wallet_single)
	}
	var err error
	taskDo := func(t tps.Task, id int) error {
		task, ok := t.(*EthTask)
		if !ok {
			return errors.New("unexpected task type")
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var (
			priv         = wallet_list[id].Priv(id)
			currentNonce = wallet_list[id].CurrentNonce(priv)
		)
		err = task.Do(ctx, &(*client)[id], priv, currentNonce, &queue, logger, TokenBankAddress)
		wallet_list[id].IncrementNonce(priv)
		wallet_list[id].IncrementNonce(priv)
		if err != nil {
			if strings.Contains(err.Error(), "ValidateInputAdvanced: Got") {
				pattern := regexp.MustCompile(`(\d+)`)
				numberStrings := pattern.FindAllStringSubmatch(err.Error(), -1)
				numbers := make([]int, len(numberStrings))
				for i, numberString := range numberStrings {
					number, err := strconv.Atoi(numberString[1])
					if err != nil {
						panic(err)
					}
					numbers[i] = number
				}
				wallet_list[id].RecetNonce(priv, uint64(numbers[3]))
				fmt.Println("Restnonce is", wallet_list[id].CurrentNonce(priv))
				return nil
			}
			if errors.Is(err, tps.ErrWrongNonce) {
				wallet_list[id].RecetNonce(priv, currentNonce)
				task.tryCount = 0
				queue.Push(task)
				return nil
			}
			if errors.Is(err, tps.ErrTaskRetry) {
				wallet_list[id].IncrementNonce(priv)
				return nil
			}
			return errors.Wrap(err, "err Do")
		}

		return nil
	}

	worker := tps.NewWorker(taskDo)

	if concurrency > MaxConcurrency {
		logger.Warn(fmt.Sprintf("concurrency setting is over logical max(%d)", MaxConcurrency))
	}
	for i := 0; i < concurrency; i++ {
		go worker.Run(&queue, i)
	}

	go func() {
		count := 2
		for {
			if atomic.LoadUint32(&closing) == 1 {
				break
			}

			if queue.CountTasks() > queueSize {
				continue
			}

			queue.Push(&EthTask{
				to:      "0x27F6F1bb3e2977c3CB014e7d4B5639bB133A6032",
				amount:  1,
				tokenId: int64(count),
			})
			count++
		}
	}()
}
