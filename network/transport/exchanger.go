package transport

import (
	"bytes"
	"context"
	"github.com/pingcap/errors"
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

func (e *exchanger) isAvailable() bool {
	select {
	case <-e.ctx.Done():
		return false
	default:
		return true
	}
}

func (e *exchanger) writePkg(messageWrappers []*MessageWrapper) map[uint32]error {
	var errMap = make(map[uint32]error)
	if messageWrappers == nil || len(messageWrappers) == 0 {
		return nil
	}
	if !e.isAvailable() {
		err := ExchangerNotAvailable
		for _, wrapper := range messageWrappers {
			errMap[wrapper.msg.PackageId] = err
		}
		return errMap
	}

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				logger.Error("error happened when writing pkg,msg:%s", errors.ErrorStack(err))
			}
		}
	}()

	if len(messageWrappers) == 1 {
		bytes, err := e.message2Byte(messageWrappers[0])
		if err != nil {
			errMap[messageWrappers[0].msg.PackageId] = err
			return errMap
		}
		if _, err := e.conn.send(bytes); err != nil {
			errMap[messageWrappers[0].msg.PackageId] = err
			return errMap
		}
	} else {
		var buffers []*bytes.Buffer
		for _, wrapper := range messageWrappers {
			bytes, err := e.message2Byte(wrapper)
			if err != nil {
				errMap[wrapper.msg.PackageId] = err
				continue
			}
			buffers = append(buffers, bytes)
		}
		if _, err := e.conn.send(buffers); err != nil {
			for _, wrapper := range messageWrappers {
				if errMap[wrapper.msg.PackageId] == nil {
					errMap[wrapper.msg.PackageId] = err
				}
			}
			return errMap
		}
	}
	return nil
}

func (e *exchanger) message2Byte(wrapper *MessageWrapper) (*bytes.Buffer, error) {
	marshalStartTime := time.Now()

	msg := wrapper.msg
	codec := wrapper.c
	if msg.Content != nil {
		if body, err := codec.Write(msg.Content); err != nil {
			return nil, err
		} else {
			if len(body) > e.compressionThreshold {
				compressed, compressError := e.compression.Compress(body)
				if compressError == nil {
					msg.Flag = compressFeature.enable(msg.Flag)
					body = compressed
				} else {
					logger.Error("compress feature msg body failed...", err)
					msg.Flag = compressFeature.disable(msg.Flag)
				}
			} else {
				msg.Flag = compressFeature.disable(msg.Flag)
			}
			msg.Body = body
		}
	}

	databytes, err := encodeMsg(msg)
	if err != nil {
		return nil, err
	} else {
		if wrapper.call != nil {
			wrapper.call.SerializeTime = time.Now().Sub(marshalStartTime)
			wrapper.call.SerializeSize = msg.Length
		}
	}
	return databytes, nil
}

func (e *exchanger) handleMsgError(packageId uint32, err error) {
	e.minusConcurrencyNum()
	call := e.removePending(packageId)
	if call == nil {
		return
	}
	call.Response = nil
	call.Error = err
	call.Done <- struct{}{}
}

func (e *exchanger) removePending(packageId uint32) *NetCall {
	if e.pendingMap == nil {
		return nil
	}
	if c, ok := e.pendingMap.Load(packageId); ok {
		e.pendingMap.Delete(packageId)
		return c.(*NetCall)
	}
	return nil
}

func (e *exchanger) minusConcurrencyNum() int32 {
	return atomic.AddInt32(&e.concurrencyNum, -1)
}

func (e *exchanger) closeGracefully() {
	e.closeOnce.Do(func() {
		e.lock.Lock()
		defer e.lock.Unlock()
		logger.Debug("exchanger start to async close!,concurrence: %d,trying to close asynchronously.", e.concurrencyNum)
		//close heart beat
		e.hb.cancel()
		//close exchanger
		e.cancel()

		go func() {
			
		}()
	})
}
