package protocol

type DiscoverType int

const (
	DISCOVER_SERVER_LEADER DiscoverType = iota
	DISCOVER_SERVER_SLAVE
	DISCOVER_CLIENT
)

func (dt DiscoverType) String() string {
	switch dt {
	case DISCOVER_SERVER_SLAVE:
		return "SERVER_SLAVE"
	case DISCOVER_SERVER_LEADER:
		return "SERVER_LEADER"
	case DISCOVER_CLIENT:
		return "CLIENT"
	default:
		return "UNKNOWN"
	}
}
