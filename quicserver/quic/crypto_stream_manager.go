package quic

import "github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"

type cryptoDataHandler interface {
	HandleMessage([]byte, protocol.EncryptionLevel) bool
}

type CryptoStreamManager struct {
	cryptoHandler  cryptoDataHandler

	initialStream   cryptoStream
	handshakeStream cryptoStream
	oneRTTStream    cryptoStream
}
