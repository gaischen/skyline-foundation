package http2

import "sync"

type transportInFlow struct {
	limit   uint32
	unacked uint32
}

type quotaPool struct {
	c     chan int
	mu    sync.Mutex
	quota int
}
