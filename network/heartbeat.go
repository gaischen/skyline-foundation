package network

import (
	"context"
	"sync/atomic"
	"time"
)

type heartbeat struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cronPeriod time.Duration
	e          *exchanger
	timeout    time.Duration
}

func newHeartbeat(e *exchanger) *heartbeat {
	h := new(heartbeat)
	h.e = e
	ctx, cancel := context.WithCancel(e.ctx)
	h.ctx, h.cancel = ctx, cancel
	h.cronPeriod = 5 * time.Second
	h.timeout = 3 * h.cronPeriod
	return h
}

func (h *heartbeat) run() {
	defer func() {
		logger.Debug("exchanger heartbeat shutdown...")
		atomic.AddInt32(&h.e.goroutineNum, -1)
	}()
	t := time.NewTicker(h.cronPeriod)
	for {
		select {
		case <-t.C:
			h.doHeartbeat()
		case <-h.ctx.Done():
			return
		}
	}
}

func (h *heartbeat) doHeartbeat() {
	if h == nil || h.e == nil || h.e.conn == nil {
		return
	}
	currentTime := time.Now().Unix()

	if (atomic.LoadInt64(&h.e.conn.readLastTime) != 0 && (currentTime-atomic.LoadInt64(&h.e.conn.readLastTime) > int64(h.cronPeriod.Seconds()))) ||
		(atomic.LoadInt64(&h.e.conn.writeLastTime) != 0 && (currentTime-atomic.LoadInt64(&h.e.conn.writeLastTime) > int64(h.cronPeriod.Seconds()))) {
		msg := new(Message)
		msg.Version = uint8(2)
		msg.Flag = 1

		// 。。。。。send

	}

}
