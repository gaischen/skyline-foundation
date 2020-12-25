package quic

import (
	"errors"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/logging"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/utils"
	"net"
)

func newPacketHandlerMap(c net.PacketConn,
	connIDLen int,
	statelessResetKey []byte,
	tracer logging.Tracer,
	logger utils.Logger) (packetHandlerManager, error) {
	if err := setReceiveBuffer(c, logger); err != nil {

	}
}

func setReceiveBuffer(c net.PacketConn, logger utils.Logger) error {
	conn, ok := c.(interface{ SetReadBuffer(int) error })
	if !ok {
		return errors.New("connection doesn't allow setting of receive buffer")
	}


}
