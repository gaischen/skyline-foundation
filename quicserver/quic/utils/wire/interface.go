package wire

import (
	"bytes"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
)

// A Frame in QUIC
type Frame interface {
	Write(b *bytes.Buffer, version protocol.VersionNumber) error
	Length(version protocol.VersionNumber) protocol.ByteCount
}