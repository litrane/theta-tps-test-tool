package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/tak1827/blockchain-tps-test/tps"
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

func (t *EthTask) Do(ctx context.Context, client *EthClient, priv string, nonce uint64, queue *tps.Queue, logger tps.Logger, erc721address string) error {

	var rootErr error

	_, rootErr = client.SendTx(ctx, priv, nonce, t.to, t.amount)

	if rootErr != nil {
		if strings.Contains(rootErr.Error(), "Invalid Transaction") {
			//logger.Warn(fmt.Sprintf("nonce error, %s", rootErr.Error()))
			return tps.ErrWrongNonce
		}
		if strings.Contains(rootErr.Error(), "Transaction already seen") {
			//logger.Warn(fmt.Sprintf("nonce error, %s", rootErr.Error()))
			return tps.ErrTaskRetry
		}
		logger.Warn(fmt.Sprintf("faild sending, err: %s", rootErr.Error()))
		if err := t.IncrementTryCount(); err != nil {
			return errors.Wrap(rootErr, err.Error())
		}
		return tps.ErrWrongNonce
		queue.Push(t)
		return nil
	}

	return nil
}
