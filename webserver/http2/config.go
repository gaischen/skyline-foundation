package http2

import "github.com/vanga-top/skyline-foundation/webserver/credentials"

type ServerConfig struct {
	MaxStream uint32
	AuthInfo  credentials.AuthInfo
}


