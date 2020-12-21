package quic

import "time"

type zeroRTTQueueEntry struct {
	timer *time.Timer

}

type zeroRTTQueue struct {
}
