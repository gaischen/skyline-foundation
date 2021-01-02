package server

import "github.com/vanga-top/skyline-foundation/discover/internal/protocol"

//处理服务端数据的定义
//goland:noinspection ALL
type ServerDataProcessor interface {
	Online(meta protocol.ServiceMeta) error
	Offline(meta protocol.ServiceMeta) error
	Register(meta protocol.ServiceMeta) error
	Remove(meta protocol.ServiceMeta) error
}

