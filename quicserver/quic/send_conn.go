package quic

import "net"

type sendConn interface {
	Write([]byte) error
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}
