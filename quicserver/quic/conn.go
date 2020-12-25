package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"io"
	"net"
	"syscall"
	"time"
)

type connection interface {
	ReadPacket() (*receivedPacket, error)
	WriteTo([]byte, net.Addr) (int, error)
	LocalAddr() net.Addr
	io.Closer
}

// If the PacketConn passed to Dial or Listen satisfies this interface, quic-go will read the ECN bits from the IP header.
// In this case, ReadMsgUDP() will be used instead of ReadFrom() to read packets.
type ECNCapablePacketConn interface {
	net.PacketConn
	SyscallConn() (syscall.RawConn, error)
	ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *net.UDPAddr, err error)
}

var _ ECNCapablePacketConn = &net.UDPConn{}

func wrapConn(pc net.PacketConn) (connection, error) {
	c, ok := pc.(ECNCapablePacketConn)
	if !ok {
		return &basicConn{pc}, nil
	}
	return newConn(c)
}

type basicConn struct {
	net.PacketConn
}

var _ connection = &basicConn{}

func (c *basicConn) ReadPacket() (*receivedPacket, error) {
	buffer := getPacketBuffer()
	buffer.Data = buffer.Data[:protocol.MaxReceivePacketSize]
	n, addr, err := c.PacketConn.ReadFrom(buffer.Data)
	if err != nil {
		return nil, err
	}
	return &receivedPacket{
		remoteAddr: addr,
		rcvTime:    time.Now(),
		data:       buffer.Data[:n],
		buffer:     buffer,
	}, nil
}
