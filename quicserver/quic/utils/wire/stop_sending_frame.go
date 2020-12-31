package wire

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

type StopSendingFrame struct {
	StreamID protocol.StreamID
	ErrorCode protocol.ApplicationErrorCode
}
