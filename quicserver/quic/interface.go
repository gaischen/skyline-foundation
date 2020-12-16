package quic

import (
	"context"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"io"
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

// A Session is a QUIC connection between two peers.
type Session interface {
	AcceptStream(ctx context.Context) (Stream, error)
	//returns the next unidirectional stream opened by the peer, blocking until one is available.
	AcceptUniStream(ctx context.Context) (ReceiveStream, error)
	
}
