package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"sync"
	"time"
)

type zeroRTTQueueEntry struct {
	timer   *time.Timer
	packets []*receivedPacket
}

type zeroRTTQueue struct {
	mutex         sync.Mutex
	queue         map[string]*zeroRTTQueueEntry
	queueDuration time.Duration
}

func newZeroRTTQueue() *zeroRTTQueue {
	return &zeroRTTQueue{
		queue:         make(map[string]*zeroRTTQueueEntry),
		queueDuration: protocol.Max0RTTQueueingDuration,
	}
}
