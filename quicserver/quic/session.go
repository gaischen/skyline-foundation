package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"net"
	"time"
)

type receivedPacket struct {
	buffer *packetBuffer

	remoteAddr net.Addr
	rcvTime    time.Time
	data       []byte
	ecn        protocol.ECN
}
