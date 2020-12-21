package quic

import (
	"crypto/tls"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/handshake"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"net"
	"sync"
)

//start listen quic addr
func ListenAddr(addr string, tlsConfig *tls.Config, config *Config) (Listener, error) {

	return nil, nil
}

type packetHandler interface {
	handlePacket(packet *receivedPacket)
	shutdown()
	destroy(error)
	getPerspective() protocol.Perspective
}

type unknownPacketHandler interface {
	handlePacket(*receivedPacket)
	setCloseError(error)
}

type packetHandlerManager interface {
	AddWithConnID(protocol.ConnectionID, protocol.ConnectionID, func() packetHandler) bool
	Destroy() error
	sessionRunner
	SetServer(unknownPacketHandler)
	CloseServer()
}

type basicServer struct {
	mutex               sync.Mutex
	acceptEarlySessions bool
	tlsConfig           *tls.Config
	config              *Config

	conn net.PacketConn
	//if the server is started with listenAddr we create a packet conn.
	//if it is started with listen we take a packet conn as a parameter
	createdPacketConn bool
	tokenGenerator    *handshake.TokenGenerator
	zeroRTTQueue      *zeroRTTQueue
	sessionHandler    packetHandlerManager
	receivePackets    chan *receivedPacket
	newSession        func(
		sendConn,
		sessionRunner,
		protocol.ConnectionID,  /* original dest connection ID */
		*protocol.ConnectionID, /* retry src connection ID */
		protocol.ConnectionID,  /* client dest connection ID */
		protocol.ConnectionID,  /* destination connection ID */
		protocol.ConnectionID,  /* source connection ID */
		protocol.StatelessResetToken,
		
	)
}
