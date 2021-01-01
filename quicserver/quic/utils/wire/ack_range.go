package wire

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

type AckRange struct {
	Smallest protocol.PacketNumber
	Largest  protocol.PacketNumber
}

// Len returns the number of packets contained in this ACK range
func (r AckRange) Len() protocol.PacketNumber {
	return r.Largest - r.Smallest + 1
}