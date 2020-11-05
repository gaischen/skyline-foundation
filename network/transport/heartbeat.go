package transport

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
	
}
