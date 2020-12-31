package wire

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

// A MaxStreamDataFrame is a MAX_STREAM_DATA frame
type MaxStreamDataFrame struct {
	StreamID          protocol.StreamID
	MaximumStreamData protocol.ByteCount
}
