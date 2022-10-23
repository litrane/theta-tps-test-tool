package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/blockchain-tps-test/samples/theta/tps"
)

const (
	Eth tps.TaskType = iota

	TaskRetryLimit = 10
)

type EthTask struct {
	to       string
	amount   int64
	tokenId  int64
	tryCount int
}

func (t *EthTask) Type() tps.TaskType {
	return Eth
}

func (t *EthTask) TryCount() int {
	return t.tryCount
}

func (t *EthTask) IncrementTryCount() error {
	t.tryCount += 1
	if t.tryCount >= TaskRetryLimit {
		return fmt.Errorf("err task retry limit, tryCount: %d", t.tryCount)
	}
	return nil
}

func (t *EthTask) Do(ctx context.Context, client *EthClient, priv string, nonce uint64, queue *tps.Queue, logger tps.Logger, contractAddress string) error {

	var rootErr error
	if model == "ERC20" {
		_, rootErr = client.Erc20TransferFrom(ctx, priv, nonce, t.to, t.amount, contractAddress, 1)
	} else if model == "CrossChainTNT20" {
		_, rootErr = client.CrossChainTNT20Transfer(ctx, priv, nonce, t.to, t.amount, contractAddress, 1)
	} else if model == "CrossSubChainTNT20" {
		_, rootErr = client.CrossSubChainTNT20Transfer(ctx, priv, nonce, t.to, t.amount, contractAddress, 1)
	} else {
		_, rootErr = client.SendTx(ctx, priv, nonce, t.to, t.amount)
	}

	if rootErr != nil {
		if strings.Contains(rootErr.Error(), "Invalid Transaction") {
			//logger.Warn(fmt.Sprintf("nonce error, %s", rootErr.Error()))
			return tps.ErrWrongNonce
		}
		if strings.Contains(rootErr.Error(), "Transaction already seen") {
			//logger.Warn(fmt.Sprintf("nonce error, %s", rootErr.Error()))
			return tps.ErrTaskRetry
		}
		if strings.Contains(rootErr.Error(), "ValidateInputAdvanced: Got") {

			return rootErr
		}
		logger.Warn(fmt.Sprintf("faild sending, err: %s", rootErr.Error()))
		if err := t.IncrementTryCount(); err != nil {
			return tps.ErrWrongNonce
		}
		// if strings.Contains(rootErr.Error(), "ValidateInputAdvanced:"){
		// 	return tps.NonceWrong
		// }
		return tps.ErrWrongNonce
		queue.Push(t)
		return nil
	}

	return nil
}
