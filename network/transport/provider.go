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
	SerializeTime   time.Duration
	DeSerializeTime time.Duration
	Timeout         time.Duration
	SerializeSize   uint32
	DeSerializeSize uint32

	Done     chan struct{}
	Error    error
	Response interface{}
	Invocation *Invocation
}

type MessageWrapper struct {
	msg  *Message
	call *NetCall
	c    codec
}
