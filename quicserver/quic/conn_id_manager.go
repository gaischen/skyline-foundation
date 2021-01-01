package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils/wire"
	mrand "math/rand"
)

type connIDManager struct {
	queue utils.NewConnectionIDList

	handshakeComplete         bool
	activeSequenceNumber      uint64
	highestRetired            uint64
	activeConnectionID        protocol.ConnectionID
	activeStatelessResetToken *protocol.StatelessResetToken

	packetsSinceLastChange uint64
	rand                   *mrand.Rand
	packetsPerConnectionID uint64

	addStatelessResetToken    func(protocol.StatelessResetToken)
	removeStatelessResetToken func(protocol.StatelessResetToken)
	queueControlFrame         func(wire.Frame)
}
