package quic

import (
	"bytes"
	"fmt"
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

type connMultiplexer struct {
	mutex                   sync.Mutex
	conns                   map[string]connManager
	newPacketHandlerManager func(net.PacketConn, int, []byte, logging.Tracer, utils.Logger) (packetHandlerManager, error) // so it can be replaced in the tests
	logger                  utils.Logger
}

func (m *connMultiplexer) AddConn(c net.PacketConn, connIDLen int, statelessResetKey []byte, tracer logging.Tracer) (packetHandlerManager, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	addr := c.LocalAddr()
	connIndex := addr.Network() + " " + addr.String()
	p, ok := m.conns[connIndex]
	if !ok {
		manager, err := m.newPacketHandlerManager(c, connIDLen, statelessResetKey, tracer, m.logger)
		if err != nil {
			return nil, err
		}
		p = connManager{
			connIDLen:         connIDLen,
			statelessResetKey: statelessResetKey,
			manager:           manager,
			tracer:            tracer,
		}
		m.conns[connIndex] = p
	} else {
		if p.connIDLen != connIDLen {
			return nil, fmt.Errorf("cannot use %d byte connection IDs on a connection that is already using %d byte connction IDs", connIDLen, p.connIDLen)
		}
		if statelessResetKey != nil && !bytes.Equal(p.statelessResetKey, statelessResetKey) {
			return nil, fmt.Errorf("cannot use different stateless reset keys on the same packet conn")
		}
		if tracer != p.tracer {
			return nil, fmt.Errorf("cannot use different tracers on the same packet conn")
		}
	}
	return p.manager, nil
}

func (m *connMultiplexer) RemoveConn(conn indexableConn) error {
	panic("implement me")
}

func getMultiplexer() multiplexer {
	connMuxerOnce.Do(func() {
		connMuxer = &connMultiplexer{
			conns:                   make(map[string]connManager),
			logger:                  utils.DefaultLogger.WithPrefix("muxer"),
			newPacketHandlerManager: newPacketHandlerMap,
		}
	})
	return connMuxer
}
