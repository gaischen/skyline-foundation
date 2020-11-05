package transport

import (
	"github.com/vanga-top/skyline-foundation/log"
	"github.com/vanga-top/skyline-foundation/log/level"
)

var logger log.Logger = log.NewLogger("transport", level.WARN)

type Packet struct {
	packetId string // channelId+auto inc

}
