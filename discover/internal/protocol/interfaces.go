package protocol

import "time"

//
type Discover interface {
}

type ServiceMeta struct {
	ServiceName    string
	ServiceVersion string
	ServiceGroup   string
	ServiceTag     []string          //tag for service
	ServiceContent map[string]string //content

	RegisterTime time.Time
	OnlineTime   time.Time
	OfflineTime  time.Time //最近一次的下线时间
}

type Service interface {
	GetServiceMeta() *ServiceMeta
	Register(sm *ServiceMeta)
	Online(duration time.Duration, timout time.Duration) error //duration==0 means publish right now
	Offline(duration time.Duration, timout time.Duration) error
}

type Changer func()

type Watcher interface {
	Watch(changer Changer)
}
