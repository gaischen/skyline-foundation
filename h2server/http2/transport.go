package http2

import (
	"github.com/vanga-top/skyline-foundation/h2server/codes"
	"github.com/vanga-top/skyline-foundation/h2server/metadata"
	"golang.org/x/net/http2"
	"net"
	"sync"
)

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

func (b *controlBuffer) put(r item) {
	b.mu.Lock()
	if len(b.backlog) == 0 {
		select {
		case b.c <- r:
			b.mu.Unlock()
			return
		default:
		}
	}
	b.backlog = append(b.backlog, r)
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
	SetId(id string)
	HandleStream(func(stream *Stream))
	WriteHeader(s *Stream, md metadata.MD) error
	Write(s *Stream, data []byte, opts *Options) error
	WriteStatus(s *Stream, statusCode codes.Code, statusDesc string) error
	Push(s *Stream, data []byte, flags http2.Flags) error
	Close() error
	RemoteAddr() net.Addr
	Drain()
}

type Options struct {
	Last  bool
	Delay bool
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
