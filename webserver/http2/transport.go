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

type transportState int

const (
	reachable transportState = iota
	unreachable
	closing
	draining
)

