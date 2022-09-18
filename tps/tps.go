package tps

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func StartTPSMeasuring(ctx context.Context, client Client, closing, idlingDuration *uint32, logger Logger) error {
	var (
		idling      = true
		startedAd   time.Time
		total       int
		count       int
		lastBlock   uint64
		err         error
		avg_latency time.Duration
	)

	for {
		if atomic.LoadUint32(closing) == 1 {
			break
		}
		if count, lastBlock, avg_latency, err = countTx(ctx, client, lastBlock); err != nil {
			if errors.Is(err, ErrNotNewBlock) {
				// sleep a bit
				//time.Sleep(1 * time.Second)
				logger.Warn("ErrNotNewBlock")
				continue
			}
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Warn("timeout of countTx")
				continue
			}
			//TODO: handle timeout error
			return errors.Wrap(err, "err CountTx")
		}

		if idling {
			if count > 0 {
				idling = false
				startedAd = time.Now()
			}
			continue
		}

		pendingTx, err := client.CountPendingTx(ctx)
		if err != nil {
			return errors.Wrap(err, "err CountPendingTx")
		}

		// NextIdlingDuration(idlingDuration, uint32(count), uint32(pendingTx))

		total += count
		elapsed := time.Now().Sub(startedAd).Seconds()
		fmt.Print("------------------------------------------------------------------------------------\n")
		fmt.Printf("⛓  %d th Block Mind! txs(%d), total txs(%d), TPS(%.2f), pendig txs(%d),latency(%d)\n", lastBlock, count, total, float64(total)/elapsed, pendingTx, avg_latency)
	}

	return nil
}

func countTx(ctx context.Context, client Client, lastBlock uint64) (int, uint64, time.Duration, error) {
	height, err := client.LatestBlockHeight(ctx)
	if err != nil {
		return 0, lastBlock, time.Duration(0), errors.Wrap(err, "err LatestBlockHeight")
	}
	if height <= lastBlock {
		return 0, lastBlock, time.Duration(0), ErrNotNewBlock
	}

	count, avg_latency, err := client.CountTx(ctx, height)
	if err != nil {
		return 0, lastBlock, time.Duration(0), errors.Wrap(err, "err TxCount")
	}

	return count, height, avg_latency, nil
}
