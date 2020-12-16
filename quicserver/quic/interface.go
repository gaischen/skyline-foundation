package quic

import (
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
	"time"
)

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


