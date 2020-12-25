package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"sync"
)

type packetBuffer struct {
	Data []byte
	// refCount counts how many packets Data is used in.
	// It doesn't support concurrent use.
	// It is > 1 when used for coalesced packet.
	refCount int
}

var bufferPool sync.Pool

func getPacketBuffer() *packetBuffer {
	buf := bufferPool.Get().(*packetBuffer)
	buf.refCount = 1
	buf.Data = buf.Data[:0]
	return buf
}

func init() {
	bufferPool.New = func() interface{} {
		return &packetBuffer{
			Data: make([]byte, 0, protocol.MaxReceivePacketSize),
		}
	}
}