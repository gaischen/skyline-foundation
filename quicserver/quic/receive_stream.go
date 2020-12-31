package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils/wire"
)

type receiveStreamI interface {
	ReceiveStream

	handleStreamFrame(*wire.StreamFrame) error
	handleResetStreamFrame(*wire.ResetStreamFrame) error
	closeForShutdown(error)
	getWindowUpdate() protocol.ByteCount
}
