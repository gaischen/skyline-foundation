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

func (b *controlBuffer) get() <-chan item {
	return b.c
}

func (b *controlBuffer) load() {
	b.mu.Lock()
	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			b.backlog[0] = nil
			b.backlog = b.backlog[1:]
		default:
		}
	}
	b.mu.Unlock()
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

type clientTransport interface {
}

type recvBuffer struct {
	c       chan recvMsg
	mu      sync.Mutex
	backlog []recvMsg
}

type recvMsg struct {
	data []byte
	err  error
}
