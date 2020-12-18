package quic

import (
	"context"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/handshake"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/logging"
	"io"
	"net"
	"time"
)

type StreamID = protocol.StreamID
type VersionNumber = protocol.VersionNumber

type Token struct {
	IsRetryToken bool
	RemoteAddr   string
	SentTime     time.Time
}

type ClientToken struct {
	data []byte
}

type TokenStore interface {
	Pop(key string) (token *ClientToken)
	Put(key string, token *ClientToken)
}

type ErrorCode = protocol.ApplicationErrorCode

type ReceiveStream interface {
	StreamID() StreamID
	io.Reader
	CancelRead(code ErrorCode)
	SetReadDeadline(t time.Time) error
}

type SendStream interface {
	StreamID() StreamID
	io.Writer
	io.Closer
	CancelWrite(code ErrorCode)
	Context() context.Context
	SetWriteDeadline(t time.Time) error
}

type Stream interface {
	ReceiveStream
	SendStream
	SetDeadline(t time.Time) error
}

type StreamError interface {
	error
	Canceled() bool
	ErrorCode() ErrorCode
}

type ConnectionState = handshake.ConnectionState

// A Session is a QUIC connection between two peers.
type Session interface {
	AcceptStream(ctx context.Context) (Stream, error)
	//returns the next unidirectional stream opened by the peer, blocking until one is available.
	AcceptUniStream(ctx context.Context) (ReceiveStream, error)
	OpenStream() (Stream, error)
	OpenStreamSync(ctx context.Context) (Stream, error)
	OpenUniStream() (SendStream, error)
	OpenUniStreamSync(ctx context.Context) (SendStream, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	CloseWithError(code ErrorCode, msg string) error
	Context() context.Context
	ConnectionState() ConnectionState
}

// An EarlySession is a session that is handshaking.
type EarlySession interface {
	Session
	HandshakeComplete() context.Context
}

type Config struct {
	Versions           []VersionNumber
	ConnectionIDLength int
	HandshakeTimeout   time.Duration
	MaxIdleTimeout     time.Duration
	AcceptToken        func(clientAddr net.Addr, token *Token) bool
	TokenStore
	// MaxReceiveStreamFlowControlWindow is the maximum stream-level flow control window for receiving data.
	// If this value is zero, it will default to 1 MB for the server and 6 MB for the client.
	NaxReceiveStreamFlowControlWindow uint64
	// MaxReceiveConnectionFlowControlWindow is the connection-level flow control window for receiving data.
	// If this value is zero, it will default to 1.5 MB for the server and 15 MB for the client.
	MaxReceiveConnectionFlowControlWindow uint64
	MaxIncomingStreams int64
	MaxIncomingUniStreams int64
	StatelessResetKey []byte
	KeepAlive bool
	Tracer    logging.Tracer
}

// A Listener for incoming QUIC connections
type Listener interface {
	// Close the server. All active sessions will be closed.
	Close() error
	// Addr returns the local network addr that the server is listening on.
	Addr() net.Addr
	// Accept returns new sessions. It should be called in a loop.
	Accept(context.Context) (Session, error)
}

// An EarlyListener listens for incoming QUIC connections,
// and returns them before the handshake completes.
type EarlyListener interface {
	// Close the server. All active sessions will be closed.
	Close() error
	// Addr returns the local network addr that the server is listening on.
	Addr() net.Addr
	// Accept returns new early sessions. It should be called in a loop.
	Accept(context.Context) (EarlySession, error)
}