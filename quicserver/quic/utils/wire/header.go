package wire

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"io"
)

func ParseConnectionID(data []byte, shortHeaderConnIDLen int) (protocol.ConnectionID, error) {
	if len(data) == 0 {
		return nil, io.EOF
	}
	isLongHeader := data[0]&0x80 > 0
	if !isLongHeader {
		if len(data) < shortHeaderConnIDLen+1 {
			return nil, io.EOF
		}
		return protocol.ConnectionID(data[1 : 1+shortHeaderConnIDLen]), nil
	}
	if len(data) < 6 {
		return nil, io.EOF
	}
	destConnIDLen := int(data[5])
	if len(data) < 6+destConnIDLen {
		return nil, io.EOF
	}
	return protocol.ConnectionID(data[6 : 6+destConnIDLen]), nil
}
