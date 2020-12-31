package wire

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

type StreamFrame struct {
	StreamID       protocol.StreamID
	Offset         protocol.ByteCount
	Data           []byte
	Fin            bool
	DataLenPresent bool

	fromPool bool
}
