package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	"context"

	"github.com/blockchain-tps-test/samples/theta/tps"
	"github.com/pkg/errors"
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
func ethStressTest(client *[]EthClient, ctx context.Context) {

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
		fmt.Println("fromAddress: ", fromAddress.Hex(), privs)
	}
	var wallet_list []tps.Wallet
	for i := 0; i < client_number; i++ {
		wallet_single, err := tps.NewWallet(ctx, (*client)[i], privs, addrs)
		if err != nil {
			logger.Fatal("err NewWallet: ", err)
		}
		wallet_list = append(wallet_list, wallet_single)
	}

	taskDo := func(t tps.Task, id int) error {
		task, ok := t.(*EthTask)
		if !ok {
			return errors.New("unexpected task type")
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var (
			priv         = wallet_list[id].Priv(id)
			currentNonce = wallet_list[id].IncrementNonce(priv)
		)
		if err := task.Do(ctx, &(*client)[id], priv, currentNonce, &queue, logger, ""); err != nil {
			if errors.Is(err, tps.ErrWrongNonce) {
				pri, _ := hex.DecodeString(priv)
				thetaPrivateKey, _ := crypto.PrivateKeyFromBytes(pri)

				nonce, err := tps.NewNonce(context.Background(), &(*client)[id], thetaPrivateKey.PublicKey().Address().Hex())
				if err != nil {
					return errors.Wrap(err, "debug")
				}
				wallet_list[id].RecetNonce(priv, nonce.Current())
				return nil
			} else if errors.Is(err, tps.ErrTaskRetry) {
				wallet_list[id].RecetNonce(priv, wallet_list[id].IncrementNonce(priv))
				return nil
			} else if errors.Is(err, tps.NonceWrong) {
				pri, _ := hex.DecodeString(priv)
				thetaPrivateKey, _ := crypto.PrivateKeyFromBytes(pri)

				nonce, err := tps.NewNonce(context.Background(), &(*client)[id], thetaPrivateKey.PublicKey().Address().Hex())
				if err != nil {
					return errors.Wrap(err, "debug")
				}
				wallet_list[id].RecetNonce(priv, nonce.Current())
				return nil
			}
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

	go worker.Run(&queue, clientID)


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
