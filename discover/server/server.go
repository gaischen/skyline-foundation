package server

import (
	"github.com/vanga-top/skyline-foundation/discover/internal/config"
	"github.com/vanga-top/skyline-foundation/discover/internal/protocol"
)

type Server interface {
	protocol.Discover
	Listen(addr string, port int) (Server, error)
	Start(serviceConfig *config.ServerConfig) (Server, error)
	Restart(serviceConfig *config.ServerConfig) (Server, error)
	Shutdown(gracefully bool) error
}

type basicServer struct {
	serverID     string
	addr         string
	port         int
	discoverType protocol.DiscoverType

	dataProcessor  *ServerDataProcessor
	leaderSelector *ServerLeaderSelector
}
