package wire

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"net"
	"time"
)

type TransportParameters struct {
	InitialMaxStreamDataBidiLocal  protocol.ByteCount
	InitialMaxStreamDataBidiRemote protocol.ByteCount
	InitialMaxStreamDataUni        protocol.ByteCount
	InitialMaxData                 protocol.ByteCount

	MaxAckDelay      time.Duration
	AckDelayExponent uint8

	DisableActiveMigration bool

	MaxUDPPayloadSize protocol.ByteCount

	MaxUniStreamNum  protocol.StreamNum
	MaxBidiStreamNum protocol.StreamNum

	MaxIdleTimeout time.Duration

	PreferredAddress *PreferredAddress

	OriginalDestinationConnectionID protocol.ConnectionID
	InitialSourceConnectionID       protocol.ConnectionID
	RetrySourceConnectionID         *protocol.ConnectionID // use a pointer here to distinguish zero-length connection IDs from missing transport parameters

	StatelessResetToken     *protocol.StatelessResetToken
	ActiveConnectionIDLimit uint64
}


// PreferredAddress is the value encoding in the preferred_address transport parameter
type PreferredAddress struct {
	IPv4                net.IP
	IPv4Port            uint16
	IPv6                net.IP
	IPv6Port            uint16
	ConnectionID        protocol.ConnectionID
	StatelessResetToken protocol.StatelessResetToken
}