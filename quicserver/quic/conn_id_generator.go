package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils/wire"
)

type connIDGenerator struct {
	connIDLen  int
	highestSqe uint64

	activeSrcConnIDs        map[uint64]protocol.ConnectionID
	initialClientDestConnID protocol.ConnectionID

	addConnectionID        func(protocol.ConnectionID)
	getStatelessResetToken func(protocol.ConnectionID) protocol.StatelessResetToken
	removeConnectionID     func(protocol.ConnectionID)
	retireConnectionID     func(protocol.ConnectionID)
	replaceWithClosed      func(protocol.ConnectionID, packetHandler)
	queueControlFrame      func(wire.Frame)

	version protocol.VersionNumber
}
