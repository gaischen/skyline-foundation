package server

import (
	"fmt"
	"github.com/vanga-top/skyline-foundation/discover/internal/config"
	"github.com/vanga-top/skyline-foundation/discover/internal/protocol"
	"net"
	"sync"
)

type Server interface {
	protocol.Discover
	protocol.LeaderSelector
	protocol.ServerDataProcessor
	Listen(network string, addr string) Server
	Start() (Server, error)
	Restart() (Server, error)
	Shutdown(gracefully bool) error
	GetPartner() []Server
	startHeartbeat()
}

type basicServer struct {
	mutex sync.Mutex
	ln    net.Listener
	wg    sync.WaitGroup

	connChanel map[string]net.Conn //key connID

	conf         *config.ServerConfig
	serverID     string
	addr         string
	network      string
	discoverType protocol.DiscoverType
}

func (b *basicServer) Online(meta protocol.ServiceMeta) error {
	panic("implement me")
}

func (b *basicServer) Offline(meta protocol.ServiceMeta) error {
	panic("implement me")
}

func (b *basicServer) Register(meta protocol.ServiceMeta) error {
	panic("implement me")
}

func (b *basicServer) Remove(meta protocol.ServiceMeta) error {
	panic("implement me")
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
		wg:           sync.WaitGroup{},
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
	if err != nil {
		fmt.Println(err)
		panic("error in listen...")
	}
	b.ln = ln
	return b
}

func (b *basicServer) Start() (Server, error) {
	b.wg.Add(1)
	go b.startHeartbeat()
	go func() {
		for {
			conn, err := b.ln.Accept()
			if err != nil {
				continue
			}
			//conn.Read()
		}
	}()

	return b, nil
}

func (b *basicServer) Restart() (Server, error) {
	panic("implement me")
}

func (b *basicServer) Shutdown(gracefully bool) error {
	panic("implement me")
}
