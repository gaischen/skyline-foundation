package http2

import (
	"sync"
	"time"
)

const (
	//默认的滑动窗口大小
	defaultWindowSize = 65535
	//初始化滑动窗口大小
	initialWindowSize             = defaultWindowSize
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

