package wire

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

type ResetStreamFrame struct {
	StreamID  protocol.StreamID
	ErrorCode protocol.ApplicationErrorCode
	FinalSize protocol.ByteCount
}
