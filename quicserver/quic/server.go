package quic

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/handshake"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/logging"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils"
	"net"
	"sync"
	"time"
)

var defaultAcceptToken = func(clientAddr net.Addr, token *Token) bool {
	if token == nil {
		return false
	}
	validity := protocol.TokenValidity
	if token.IsRetryToken {
		validity = protocol.RetryTokenValidity
	}

	if time.Now().After(token.SentTime.Add(validity)) {
		return false
	}
	var sourceAddr string
	if udpAddr, ok := clientAddr.(*net.UDPAddr); ok {
		sourceAddr = udpAddr.IP.String()
	} else {
		sourceAddr = clientAddr.String()
	}
	return sourceAddr == token.RemoteAddr
}

//start listen quic addr
func ListenAddr(addr string, tlsConfig *tls.Config, config *Config) (Listener, error) {
	return listenAddr(addr, tlsConfig, config, false)
}

func listenAddr(addr string, tlsConf *tls.Config, config *Config, acceptEarly bool) (*basicServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	serv, err := listen(conn, tlsConf, config, acceptEarly)
	if err != nil {
		return nil, err
	}
	serv.createdPacketConn = true
	return serv, nil
}

func listen(conn *net.UDPConn, tlsConf *tls.Config, config *Config, early bool) (*basicServer, error) {
	if nil == tlsConf {
		return nil, errors.New("quic: tls.Config is not set")
	}
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	config = populateServerConfig(config)
	for _, v := range config.Versions {
		if !protocol.IsValidVersion(v) {
			return nil, fmt.Errorf("%s is not a valid quic version", v)
		}
	}
	sessionHandler, err := getMultiplexer().AddConn(conn, config.ConnectionIDLength, config.StatelessResetKey, config.Tracer)
	if err != nil {
		return nil, err
	}
	tokenGenerator, err := handshake.NewTokenGenerator(rand.Reader)
	if err != nil {
		return nil, err
	}
	s := &basicServer{
		conn:           conn,
		tlsConfig:      tlsConf,
		config:         config,
		tokenGenerator: tokenGenerator,
		sessionHandler: sessionHandler,
		zeroRTTQueue:   newZeroRTTQueue(),
		sessionQueue:   make(chan quicSession),
		errorChan:      make(chan struct{}),
		running:        make(chan struct{}),
		receivePackets: make(chan *receivedPacket, protocol.MaxServerUnprocessedPackets),
		newSession:     newSession,
	}

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

type quicSession interface {
	EarlySession
	earlySessionReady() <-chan struct{}
	handlePacket(*receivedPacket)
	GetVersion() protocol.VersionNumber
	getPerspective() protocol.Perspective
	run() error
	destroy(error)
	shutdown()
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
		*Config,
		*tls.Config,
		*handshake.TokenGenerator,
		bool, /* enable 0-RTT */
		logging.ConnectionTracer,
		utils.Logger,
		protocol.VersionNumber,
	) quicSession

	serverError error
	errorChan   chan struct{}
	closed      bool
	running     chan struct{}

	sessionQueue    chan quicSession
	sessionQueueLen int32 // to be used as an atomic

	logger utils.Logger
}

func (b *basicServer) Close() error {
	panic("implement me")
}

func (b *basicServer) Addr() net.Addr {
	panic("implement me")
}

func (b *basicServer) Accept(ctx context.Context) (Session, error) {
	panic("implement me")
}
