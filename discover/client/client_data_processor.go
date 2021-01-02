package client

import (
	"github.com/vanga-top/skyline-foundation/discover/internal/protocol"
	"time"
)

type ClientDataProcessor interface {
	GetServiceMeta() *protocol.ServiceMeta
	Register(sm *protocol.ServiceMeta)
	Online(duration time.Duration, timout time.Duration) error //duration==0 means publish right now
	Offline(duration time.Duration, timout time.Duration) error
	Remove() error
}
