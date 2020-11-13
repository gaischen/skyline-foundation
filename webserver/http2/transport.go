package http2

import "sync"

type item interface {
	item()
}

type controlBuffer struct {
	c       chan item
	mu      sync.Mutex
	backlog []item
}

func newControlBuffer() *controlBuffer {
	return &controlBuffer{
		c: make(chan item, 1),
	}
}

type transportState int

const (
	reachable transportState = iota
	unreachable
	closing
	draining
)

type ServerTransport interface {
}
