package transport

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type ExchangerCallback = func(exchanger *exchanger) error

const defaultWqLen int64 = 1024

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

	hb       *heartbeat
	hbPeriod time.Duration

	compression          compression
	compressionThreshold int

	pendingMap *sync.Map

	closeDone chan struct{}

	seqGen uint32
}

func newExchanger(ctx context.Context, conn *connection, id int32) (*exchanger, error) {
	sctx, cancel := context.WithCancel(ctx)
	s := &exchanger{
		id:        id,
		conn:      conn,
		ctx:       sctx,
		cancel:    cancel,
		closeDone: make(chan struct{}, 1),
		seqGen:    0,
		wcLen:     defaultWqLen,
	}
	return s, nil
}

func (e *exchanger) run() error {

	if e.writeChan == nil {
		e.writeChan = make(chan *MessageWrapper, defaultWqLen)
	}

	if e.pendingMap == nil {
		e.pendingMap = new(sync.Map)
	}

	if e.compression == nil {
		e.compression = &lz4Compression{}
	}

	e.hb = newHeartbeat(e)
	go e.hb.run()

	atomic.AddInt32(&e.goroutineNum, 2)

	//start

	return nil
}

func (e *exchanger) loop() {
	
}
