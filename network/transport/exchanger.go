package transport

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type ExchangerCallback = func(exchanger *exchanger) error

const defaultWqLen int64 = 1024
const maxIovecNum int = 10

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
	go e.loop()

	return nil
}

func (e *exchanger) loop() {
	defer func() {
		if r := recover(); r != nil {
			atomic.AddInt32(&e.goroutineNum, -1)
		}
	}()

	var iovec = make([]*MessageWrapper, maxIovecNum)

	for {
		select {
		case data, ok := <-e.writeChan:
			if !ok {
				continue
			}
			iovec = iovec[:0]
			iovec = append(iovec, data)
		LOOP:
			for i := 0; i < maxIovecNum-1; i++ {
				select {
				case data, ok = <-e.writeChan:
					if !ok {
						break
					}
					iovec = append(iovec, data)
				default:
					break LOOP
				}
			}
			errMap := e.writePkg(iovec)
			if errMap != nil {
				for pkgId, err := range errMap {
					e.handleMsgError(pkgId, err)
				}
			}
		case <-e.ctx.Done():
			if len(e.writeChan) != 0 {
				continue
			}
			e.closeGracefully()
			return
		}
	}
}

func (e *exchanger) writePkg(messageWrappers []*MessageWrapper) map[uint32]error {

	return nil
}

func (e *exchanger) handleMsgError(packageId uint32, err error) {}

func (e *exchanger) closeGracefully() {

}
