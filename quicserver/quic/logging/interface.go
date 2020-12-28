package logging

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"net"
)

type Tracer interface {

	DroppedPacket(net.Addr, PacketType, ByteCount, PacketDropReason)
}

type ConnectionTracer interface {

}


type (
	ByteCount = protocol.ByteCount

)