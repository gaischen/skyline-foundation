package quic

import (
	"context"
	"crypto/tls"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/ackhandler"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/handshake"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/logging"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils/wire"
	"net"
	"time"
)

type receivedPacket struct {
	buffer *packetBuffer

	remoteAddr net.Addr
	rcvTime    time.Time
	data       []byte
	ecn        protocol.ECN
}

func (p *receivedPacket) Size() protocol.ByteCount {
	return protocol.ByteCount(len(p.data))
}

type sessionRunner interface {
	Add(protocol.ConnectionID, packetHandler) bool
	GetStatelessResetToken(protocol.ConnectionID) protocol.StatelessResetToken
	Retry(protocol.ConnectionID)
	Remove(protocol.ConnectionID)
	ReplaceWithClosed(protocol.ConnectionID, packetHandler)
	AddResetToken(protocol.StatelessResetToken, packetHandler)
	RemoveResetToken(protocol.StatelessResetToken)
}

type streamManager interface {
	GetOrOpenSendStream(protocol.StreamID) (sendStreamI, error)
	GetOrOpenReceiveStream(protocol.StreamID) (receiveStreamI, error)
	OpenStream() (Stream, error)
	OpenUniStream() (SendStream, error)
	OpenStreamSync(context.Context) (Stream, error)
	OpenUniStreamSync(context.Context) (SendStream, error)
	AcceptStream(context.Context) (Stream, error)
	AcceptUniStream(context.Context) (ReceiveStream, error)
	DeleteStream(protocol.StreamID) error
	UpdateLimits(*wire.TransportParameters) error
	HandleMaxStreamsFrame(*wire.MaxStreamsFrame) error
	CloseWithError(error)
}

type session struct {
	// Destination connection ID used during the handshake.
	// Used to check source connection ID on incoming packets.
	handshakeDestConnID protocol.ConnectionID
	// Set for the client. Destination connection ID used on the first Initial sent.
	origDestConnID protocol.ConnectionID
	retrySrcConnID *protocol.ConnectionID // only set for the client (and if a Retry was performed)

	srcConnIDLen int

	perspective    protocol.Perspective
	initialVersion protocol.VersionNumber // if version negotiation is performed, this is the version we initially tried
	version        protocol.VersionNumber
	config         *Config

	conn      sendConn
	sendQueue *sendQueue

	streamsMap      streamManager
	connIDManager   *connIDManager
	connIDGenerator *connIDGenerator

	rttStats *utils.RTTStats

	cryptoStreamManager   *CryptoStreamManager
	sentPacketHandler     ackhandler.SentPacketHandler
	receivedPacketHandler ackhandler.ReceivedPacketHandler
	retransmissionQueue   *retransmissionQueue
}

func (s session) AcceptStream(ctx context.Context) (Stream, error) {
	panic("implement me")
}

func (s session) AcceptUniStream(ctx context.Context) (ReceiveStream, error) {
	panic("implement me")
}

func (s session) OpenStream() (Stream, error) {
	panic("implement me")
}

func (s session) OpenStreamSync(ctx context.Context) (Stream, error) {
	panic("implement me")
}

func (s session) OpenUniStream() (SendStream, error) {
	panic("implement me")
}

func (s session) OpenUniStreamSync(ctx context.Context) (SendStream, error) {
	panic("implement me")
}

func (s session) LocalAddr() net.Addr {
	panic("implement me")
}

func (s session) RemoteAddr() net.Addr {
	panic("implement me")
}

func (s session) CloseWithError(code ErrorCode, msg string) error {
	panic("implement me")
}

func (s session) Context() context.Context {
	panic("implement me")
}

func (s session) ConnectionState() ConnectionState {
	panic("implement me")
}

func (s session) HandshakeComplete() context.Context {
	panic("implement me")
}

func (s session) earlySessionReady() <-chan struct{} {
	panic("implement me")
}

func (s session) handlePacket(packet *receivedPacket) {
	panic("implement me")
}

func (s session) GetVersion() protocol.VersionNumber {
	panic("implement me")
}

func (s session) getPerspective() protocol.Perspective {
	panic("implement me")
}

func (s session) run() error {
	panic("implement me")
}

func (s session) destroy(err error) {
	panic("implement me")
}

func (s session) shutdown() {
	panic("implement me")
}

var newSession = func(
	conn sendConn,
	runner sessionRunner,
	origDestConnID protocol.ConnectionID,
	retrySrcConnID *protocol.ConnectionID,
	clientDestConnID protocol.ConnectionID,
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	statelessResetToken protocol.ConnectionID,
	conf *Config,
	tlsConfig *tls.Config,
	tokenGenerator *handshake.TokenGenerator,
	enable0RTT bool,
	log logging.ConnectionTracer,
	tracer logging.Tracer,
	v protocol.VersionNumber,
) quicSession {
	s := &session{
		conn:                  conn,
		config:                conf,
		handshakeDestConnID:   destConnID,
		srcConnIDLen:          srcConnID.Len(),
		tokenGenerator:        tokenGenerator,
		oneRTTStream:          newCryptoStream(),
		perspective:           protocol.PerspectiveServer,
		handshakeCompleteChan: make(chan struct{}),
		tracer:                tracer,
		logger:                logger,
		version:               v,
	}

	return s
}
