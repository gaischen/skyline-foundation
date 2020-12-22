package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/logging"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils"
	"net"
	"sync"
)

var (
	connMuxerOnce sync.Once
	connMuxer     multiplexer
)

type indexableConn interface {
	LocalAddr() net.Addr
}

type multiplexer interface {
	AddConn(c net.PacketConn, connIDLen int, statelessResetKey []byte, tracer logging.Tracer) (packetHandlerManager, error)
	RemoveConn(indexableConn) error
}

type connManager struct {
	connIDLen         int
	statelessResetKey []byte
	tracer            logging.Tracer
	manager           packetHandlerManager
}

type ConnMultiplexer struct {
	mutex                   sync.Mutex
	conns                   map[string]connManager
	newPacketHandlerManager func(net.PacketConn, int, []byte, logging.Tracer, utils.Logger) (packetHandlerManager, error) // so it can be replaced in the tests
	logger                  utils.Logger
}

func (c2 *ConnMultiplexer) AddConn(c net.PacketConn, connIDLen int, statelessResetKey []byte, tracer logging.Tracer) (packetHandlerManager, error) {
	panic("implement me")
}

func (c2 *ConnMultiplexer) RemoveConn(conn indexableConn) error {
	panic("implement me")
}

func getMultiplexer() multiplexer {
	connMuxerOnce.Do(func() {
		connMuxer = &ConnMultiplexer{
			conns: make(map[string]connManager),
			logger: utils.DefaultLogger.WithPrefix("muxer"),


		}
	})
}
