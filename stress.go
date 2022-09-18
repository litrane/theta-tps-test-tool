package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"sync/atomic"

	"context"

	"github.com/pkg/errors"
	"github.com/tak1827/blockchain-tps-test/tps"
	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/crypto"
	"github.com/thetatoken/theta/crypto/sha3"
)

func keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func pubkeyToAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := crypto.FromECDSAPub(&p)
	return common.BytesToAddress(keccak256(pubBytes[1:])[12:])
}
func ethStressTest(client *EthClient, ctx context.Context) {

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

	wallet, err := tps.NewWallet(ctx, client, privs, addrs)
	if err != nil {
		logger.Fatal("err NewWallet: ", err)
	}

	taskDo := func(t tps.Task, id int) error {
		task, ok := t.(*EthTask)
		if !ok {
			return errors.New("unexpected task type")
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var (
			priv         = wallet.Priv(id)
			currentNonce = wallet.IncrementNonce(priv)
		)
		if err = task.Do(ctx, client, priv, currentNonce, &queue, logger, ""); err != nil {
			if errors.Is(err, tps.ErrWrongNonce) {

				pri, _ := hex.DecodeString(priv)
				thetaPrivateKey, _ := crypto.PrivateKeyFromBytes(pri)

				nonce, err := tps.NewNonce(context.Background(), client, thetaPrivateKey.PublicKey().Address().Hex())
				if err != nil {
					return errors.Wrap(err, "debug")
				}
				fmt.Println("need", nonce)
				fmt.Println(wallet.CurrentNonce(priv))
				wallet.RecetNonce(priv, nonce.Current())
				fmt.Println(wallet.CurrentNonce(priv))
				task.tryCount = 0
				queue.Push(task)
				return nil
			} else if errors.Is(err, tps.ErrTaskRetry) {
				wallet.RecetNonce(priv, wallet.CurrentNonce(priv))
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
		count := 0
		for {
			if atomic.LoadUint32(&closing) == 1 {
				break
			}

			if queue.CountTasks() > queueSize {
				continue
			}

			queue.Push(&EthTask{
				to:     " 0x2E833968E5bB786Ae419c4d13189fB081Cc43bab",
				amount: 1, //设置打多少币 0.001
			})
			count++
		}
	}()
}
