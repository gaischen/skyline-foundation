package quic

import (
	"errors"
	"fmt"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils"
	"golang.org/x/sys/unix"
	"net"
	"syscall"
	"time"
)

const ecnMask uint8 = 0x3

func inspectReadBuffer(c net.PacketConn) (int, error) {
	conn, ok := c.(interface {
		SyscallConn() (syscall.RawConn, error)
	})
	if !ok {
		return 0, errors.New("doesn't have a SyscallConn")
	}
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return 0, fmt.Errorf("couldn't get syscall.RawConn:#{err}")
	}
	var size int
	var serr error
	if err := rawConn.Control(func(fd uintptr) {
		size, serr = unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_RCVBUF)
	}); err != nil {
		return 0, err
	}
	return size, serr
}

type ecnConn struct {
	ECNCapablePacketConn
	oobBuffer []byte
}

func (c *ecnConn) ReadPacket() (*receivedPacket, error) {
	buffer := getPacketBuffer()
	buffer.Data = buffer.Data[:protocol.MaxReceivePacketSize]
	c.oobBuffer = c.oobBuffer[:cap(c.oobBuffer)]
	n, oobn, _, addr, err := c.ECNCapablePacketConn.ReadMsgUDP(buffer.Data, c.oobBuffer)
	if err != nil {
		return nil, err
	}
	ctrlMsgs, err := unix.ParseSocketControlMessage(c.oobBuffer[:oobn])
	if err != nil {
		return nil, err
	}
	var ecn protocol.ECN
	for _, ctrlMsg := range ctrlMsgs {
		if ctrlMsg.Header.Level == unix.IPPROTO_IP && ctrlMsg.Header.Type == msgTypeIPTOS {
			ecn = protocol.ECN(ctrlMsg.Data[0] & ecnMask)
			break
		}
		if ctrlMsg.Header.Level == unix.IPPROTO_IPV6 && ctrlMsg.Header.Type == unix.IPV6_TCLASS {
			ecn = protocol.ECN(ctrlMsg.Data[0] & ecnMask)
			break
		}
	}
	return &receivedPacket{
		remoteAddr: addr,
		rcvTime:    time.Now(),
		data:       buffer.Data[:n],
		ecn:        ecn,
		buffer:     buffer,
	}, nil
}

func newConn(c ECNCapablePacketConn) (*ecnConn, error) {
	rawConn, err := c.SyscallConn()
	if err != nil {
		return nil, err
	}
	var errIPv4, errIPv6 error
	if err := rawConn.Control(func(fd uintptr) {
		errIPv4 = unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_RECVTOS, 1)
	}); err != nil {
		return nil, err
	}
	if err := rawConn.Control(func(fd uintptr) {
		errIPv6 = unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_RECVTCLASS, 1)
	}); err != nil {
		return nil, err
	}
	switch {
	case errIPv4 == nil && errIPv6 == nil:
		utils.DefaultLogger.Debugf("Activation reading of ECN bits for IPv4 and IPv6")
	case errIPv4 == nil && errIPv6 != nil:
		utils.DefaultLogger.Debugf("Activating reading of ECN bits for IPv4.")
	case errIPv4 != nil && errIPv6 == nil:
		utils.DefaultLogger.Debugf("Activating reading of ECN bits for IPv6.")
	case errIPv4 != nil && errIPv6 != nil:
		return nil, errors.New("activating ECN failed for both IPv4 and IPv6")
	}
	return &ecnConn{
		ECNCapablePacketConn: c,
		oobBuffer:            make([]byte, 128),
	}, nil
}
