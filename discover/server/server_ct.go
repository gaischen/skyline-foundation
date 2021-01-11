package server

import "github.com/vanga-top/skyline-foundation/discover/internal/protocol"

//client for server
type ServerCT struct {
	Leader              Server
	Partners            []Server
	CurrentDiscoverType protocol.DiscoverType //当前节点的type 如果是leader则不需要互联
	Port                string
	Ipv4                string
}
