package quic

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/logging"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils"
	"hash"
	"log"
	"net"
	"sync"
	"time"
)

type packetHandlerMap struct {
	mutex       sync.Mutex
	conn        connection
	connIDLen   int
	handlers    map[string]packetHandler //key connectionID
	resetTokens map[protocol.StatelessResetToken]packetHandler
	server      unknownPacketHandler

	listening chan struct{}
	closed    bool

	deleteRetriedSessionsAfter time.Duration

	statelessResetEnabled bool
	statelessResetMutex   sync.Mutex
	statelessResetHasher  hash.Hash

	tracer logging.Tracer
	logger utils.Logger
}

func (h *packetHandlerMap) AddWithConnID(id protocol.ConnectionID, id2 protocol.ConnectionID, f func() packetHandler) bool {
	panic("implement me")
}

func (h *packetHandlerMap) Destroy() error {
	panic("implement me")
}

func (h *packetHandlerMap) Add(id protocol.ConnectionID, handler packetHandler) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.handlers[string(id)]; ok {
		h.logger.Debugf("Not adding connection ID %s, as it already exists.", id)
		return false
	}
}

func (h *packetHandlerMap) GetStatelessResetToken(id protocol.ConnectionID) protocol.StatelessResetToken {
	panic("implement me")
}

func (h *packetHandlerMap) Retry(id protocol.ConnectionID) {
	panic("implement me")
}

func (h *packetHandlerMap) Remove(id protocol.ConnectionID) {
	panic("implement me")
}

func (h *packetHandlerMap) ReplaceWithClosed(id protocol.ConnectionID, handler packetHandler) {
	panic("implement me")
}

func (h *packetHandlerMap) AddResetToken(token protocol.StatelessResetToken, handler packetHandler) {
	panic("implement me")
}

func (h *packetHandlerMap) RemoveResetToken(token protocol.StatelessResetToken) {
	panic("implement me")
}

func (h *packetHandlerMap) SetServer(handler unknownPacketHandler) {
	panic("implement me")
}

func (h *packetHandlerMap) CloseServer() {
	panic("implement me")
}

func (h *packetHandlerMap) listen() {

}

func (h *packetHandlerMap) logUsage() {

}

var _ packetHandlerManager = &packetHandlerMap{}

// only print warnings about the UPD receive buffer size once
var receiveBufferWarningOnce sync.Once

func newPacketHandlerMap(c net.PacketConn,
	connIDLen int,
	statelessResetKey []byte,
	tracer logging.Tracer,
	logger utils.Logger) (packetHandlerManager, error) {
	if err := setReceiveBuffer(c, logger); err != nil {
		receiveBufferWarningOnce.Do(func() {
			log.Printf("%s. See https://github.com/lucas-clemente/quic-go/wiki/UDP-Receive-Buffer-Size for details.", err)
		})
	}
	conn, err := wrapConn(c)
	if err != nil {
		return nil, err
	}
	m := &packetHandlerMap{
		conn:                       conn,
		connIDLen:                  connIDLen,
		listening:                  make(chan struct{}),
		handlers:                   make(map[string]packetHandler),
		resetTokens:                make(map[protocol.StatelessResetToken]packetHandler),
		deleteRetriedSessionsAfter: protocol.RetiredConnectionIDDeleteTimeout,
		statelessResetEnabled:      len(statelessResetKey) > 0,
		statelessResetHasher:       hmac.New(sha256.New, statelessResetKey),
		tracer:                     tracer,
		logger:                     logger,
	}
	go m.listen()
	if logger.Debug() {
		go m.logUsage()
	}
	return m, nil
}

func setReceiveBuffer(c net.PacketConn, logger utils.Logger) error {
	conn, ok := c.(interface{ SetReadBuffer(int) error })
	if !ok {
		return errors.New("connection doesn't allow setting of receive buffer")
	}
	size, err := inspectReadBuffer(c)
	if err != nil {
		return fmt.Errorf("failed to determine receive buffer size: %w", err)
	}
	if size > protocol.DesiredReceiveBufferSize {
		logger.Debugf("Conn has receive buffer of %d kiB (wanted: at least %d kiB)", size/1024, protocol.DesiredReceiveBufferSize/1024)
	}
	if err := conn.SetReadBuffer(protocol.DesiredReceiveBufferSize); err != nil {
		return fmt.Errorf("failed to increase receive buffer size: %w", err)
	}
	newSize, err := inspectReadBuffer(c)
	if err != nil {
		return fmt.Errorf("failed to determine receive buffer size: %w", err)
	}
	if newSize == size {
		return fmt.Errorf("failed to determine receive buffer size: %w", err)
	}
	if newSize < protocol.DesiredReceiveBufferSize {
		return fmt.Errorf("failed to sufficiently increase receive buffer size (was: %d kiB, wanted: %d kiB, got: %d kiB)", size/1024, protocol.DesiredReceiveBufferSize/1024, newSize/1024)
	}
	logger.Debugf("Increased receive buffer size to %d kiB", newSize/1024)
	return nil
}
