package transport

import (
	"context"
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


