package wire

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

type CryptoFrame struct {
	Offset protocol.ByteCount
	Data   []byte
}
