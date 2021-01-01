package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils/wire"
	"io"
)

type cryptoStream interface {
	// for receiving data
	HandleCryptoFrame(*wire.CryptoFrame) error
	GetCryptoData() []byte
	Finish() error
	// for sending data
	io.Writer
	HasData() bool
	PopCryptoFrame(protocol.ByteCount) *wire.CryptoFrame
}
