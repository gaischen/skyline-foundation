package server

import (
	"context"
	"fmt"
	"github.com/vanga-top/skyline-foundation/discover/internal/config"
	"github.com/vanga-top/skyline-foundation/discover/internal/protocol"
	"net"
	"os"
	"sync"
)

type Server interface {
	protocol.Discover
	protocol.LeaderSelector
	protocol.ServerDataProcessor
	Listen(network string, addr string) Server
	Start() Server
	Restart() (Server, error)
	Shutdown(gracefully bool) error
	GetPartner() []Server
	startHeartbeat()
}

type serverStatus int8

const (
	starting serverStatus = iota
	running
	stopping
)

type basicServer struct {
	mutex sync.Mutex
	ln    net.Listener
	wg    sync.WaitGroup

	connChanel     map[string]net.Conn //key connID 存储client过来的链接
	partnerChannel map[string]Server   // key serverID
	status         chan serverStatus
	ctx            context.Context
	cancelFunc     context.CancelFunc

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
	//panic("implement me")
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
	ctx, cancelFunc := context.WithCancel(context.Background())

	s := &basicServer{
		discoverType: discoveryType,
		wg:           sync.WaitGroup{},
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		status:       make(chan serverStatus),
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

func (b *basicServer) Start() Server {
	b.wg.Add(1)
	go b.startHeartbeat()
	go func() {
		for {
			if c := recover(); c != nil {
				fmt.Println(c)
				return
			}
			select {
			case s := <-b.status:
				if s == stopping {
					return
				}
			default:
				conn, err := b.ln.Accept()
				if err != nil {
					continue
				}
				go b.handleConn(conn)
			}
		}
	}()
	return b
}

func (b *basicServer) handleConn(conn net.Conn) {
	select {
	case <-b.ctx.Done():
		conn.Close()
		fmt.Println("ctx done...")
		return
	case s := <-b.status:
		if s == stopping {
			fmt.Println("get stat stopping...")
			b.cancelFunc()
		}
	default:
		defer conn.Close()
		br := make([]byte, 1024)
		ln, err := conn.Read(br)
		if err != nil {
			fmt.Println(err)
			err = conn.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println("receive packet length:", ln)
		fmt.Println("receive packet msg:", string(br))
		//conn.Write([]byte("hello..."))

		return
	}
}

func (b *basicServer) Restart() (Server, error) {
	panic("implement me")
}

func (b *basicServer) Shutdown(gracefully bool) error {
	if gracefully {
		b.status <- stopping
		b.wg.Done()
		return nil
	}
	os.Exit(0)
	return nil
}
