package protocol

import "time"

//
type Discover interface {
	DiscoverType() DiscoverType
	ID() string
}

type ServiceMeta struct {
	ServiceID      string //全局唯一
	ServiceName    string
	ServiceVersion string
	ServiceGroup   string
	ServiceTag     []string          //tag for service
	ServiceContent map[string]string //content

	RegisterTime time.Time
	OnlineTime   time.Time
	OfflineTime  time.Time //最近一次的下线时间
}

//数据变更的回掉信息
type Changer func(meta ServiceMeta, discover Discover)

type Watcher interface {
	Watch(changer Changer)
}
