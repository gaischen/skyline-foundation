package quic

import (
	"crypto/tls"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/handshake"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/logging"
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
}
