package server

import (
	"github.com/vanga-top/skyline-foundation/discover/internal/protocol"
	"net"
)

//client for server
type ServerCT struct {
	Leader              Server
	CurrentDiscoverType protocol.DiscoverType //当前节点的type 如果是leader则不需要互联
	Port                string
	Ipv4                string
}

func NewServerCT(leader Server) *ServerCT {
	addr, isLe := isLeader(leader)
	var disType protocol.DiscoverType
	if isLe {
		disType = protocol.DISCOVER_SERVER_LEADER
	} else {
		disType = protocol.DISCOVER_SERVER_SLAVE
	}
	return &ServerCT{
		Leader:              leader,
		CurrentDiscoverType: disType,
		Ipv4:                addr,
	}
}

func isLeader(leader Server) (string, bool) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic("error in get ips..")
	}
	var addr string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && ipnet.IP.To4().String() == leader.GetAddr() {
				addr = ipnet.IP.String()
				return addr, true
			} else if ipnet.IP.To4() != nil {
				addr = ipnet.IP.String()
			}
		}
	}
	return addr, false
}
