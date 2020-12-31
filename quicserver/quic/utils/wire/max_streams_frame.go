package wire

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

type MaxStreamsFrame struct {
	Type         protocol.StreamType
	MaxStreamNum protocol.StreamNum
}
