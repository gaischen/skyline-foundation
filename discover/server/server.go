package server

import (
	"github.com/vanga-top/skyline-foundation/discover/internal/config"
	"github.com/vanga-top/skyline-foundation/discover/internal/protocol"
)

type Server interface {
	Listen(addr string, port int) (Server, error)
	Start(serviceConfig *config.ServerConfig) (Server, error)
	Restart(serviceConfig *config.ServerConfig) (Server, error)
	Shutdown(gracefully bool) error
	DiscoverType() protocol.DiscoverType

}
