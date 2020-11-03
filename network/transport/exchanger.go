package transport

import (
	"context"
	"sync"
)

type ExchangerCallback = func(exchanger *exchanger) error

type exchanger struct {
	id        int32
	ctx       context.Context
	cancel    context.CancelFunc
	once      sync.Once
	closeOnce sync.Once
	lock      sync.Mutex

	conn     *connection
	provider *Provider

	writeChan chan *MessageWrapper
	wcLen     int64

	goroutineNum   int32
	concurrencyNum int32




}
