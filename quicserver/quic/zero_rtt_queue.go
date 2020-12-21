package quic

import (
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


