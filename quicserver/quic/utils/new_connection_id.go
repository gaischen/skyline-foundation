package utils

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

type NewConnectionID struct {
	SequenceNumber      uint64
	ConnectionID        protocol.ConnectionID
	StatelessResetToken protocol.StatelessResetToken
}
