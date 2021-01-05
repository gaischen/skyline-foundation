package server

import (
	"github.com/vanga-top/skyline-foundation/discover/internal/config"
	"github.com/vanga-top/skyline-foundation/discover/internal/protocol"
	"net"
	"sync"
)

type Server interface {
	protocol.Discover
	Listen(network string, addr string) Server
	Start() (Server, error)
	Restart() (Server, error)
	Shutdown(gracefully bool) error
	GetPartner() []Server
	startHeartbeat()
}

type basicServer struct {
	mutex sync.Mutex

	conf         *config.ServerConfig
	serverID     string
	addr         string
	network      string
	discoverType protocol.DiscoverType

	dataProcessor  *ServerDataProcessor
	leaderSelector *ServerLeaderSelector
}

func (b *basicServer) startHeartbeat() {
	panic("implement me")
}

func (b *basicServer) GetPartner() []Server {
	panic("implement me")
}

//
func NewBasicServer(conf *config.ServerConfig) Server {
	if conf == nil {
		conf = config.NewDefaultConfig()
	}

	discoveryType := protocol.ParseDiscoverType(conf.ServerType)

	s := &basicServer{
		discoverType: discoveryType,
	}

	return s
}

func (b *basicServer) DiscoverType() protocol.DiscoverType {
	return b.discoverType
}

func (b *basicServer) ID() string {
	return b.serverID
}

func (b *basicServer) Listen(network string, addr string) Server {
	ln, err := net.Listen(network, addr)

	return b
}

func (b *basicServer) Start() (Server, error) {
	panic("implement me")
}

func (b *basicServer) Restart() (Server, error) {
	panic("implement me")
}

func (b *basicServer) Shutdown(gracefully bool) error {
	panic("implement me")
}
