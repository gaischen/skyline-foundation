package http2

import "github.com/vanga-top/skyline-foundation/h2server/credentials"

type ServerConfig struct {
	MaxStream uint32
	AuthInfo  credentials.AuthInfo
}


