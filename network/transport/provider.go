package transport

import (
	"context"
	"sync"
	"time"
)

type ProviderOptions struct {
	Address      string
	ConnNum      int
	ConnInternal time.Duration
}

type ProviderOption func(options *ProviderOptions)

type Provider struct {
	ctx    context.Context
	cancel context.CancelFunc

	ProviderOptions
	exchangers []*exchanger
	once       sync.Once
	mutex      sync.Mutex

	callback       ExchangerCallback
	index          int64
	exchangerIndex int32

	url *URL
}

type NetCall struct {
}


type MessageWrapper struct {
	msg  *Message
	call *NetCall
	c    codec
}
