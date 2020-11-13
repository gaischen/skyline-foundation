package http2

import (
	"golang.org/x/net/http2"
	"sync"
	"time"
)

const (
	//默认的滑动窗口大小
	defaultWindowSize = 65535
	//初始化滑动窗口大小
	initialWindowSize             = defaultWindowSize
	initialConnWindowSize         = defaultWindowSize * 16
	defaultServerKeepaliveTime    = 5 * time.Minute
	defaultServerKeepaliveTimeout = 1 * time.Minute
)

type transportInFlow struct {
	limit   uint32
	unacked uint32
}

type quotaPool struct {
	c     chan int
	mu    sync.Mutex
	quota int
}

func newQuotaPool(q int) *quotaPool {
	qb := &quotaPool{
		c: make(chan int, 1),
	}
	if q > 0 {
		qb.c <- q
	} else {
		qb.quota = q
	}

	return qb
}

type windowUpdate struct {
	streamId  uint32
	increment uint32
}

func (*windowUpdate) item() {

}

type settings struct {
	ack bool
	ss  []http2.Setting
}

func (*settings) item() {

}
