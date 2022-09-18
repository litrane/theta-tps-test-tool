package tps

import (
	"context"
	"time"
)

type Client interface {
	LatestBlockHeight(context.Context) (uint64, error)
	CountTx(context.Context, uint64) (int, time.Duration,error)
	CountPendingTx(context.Context) (int, error)
	Nonce(context.Context, string) (uint64, error)
}
