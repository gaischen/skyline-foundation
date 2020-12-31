package ackhandler

import "github.com/vanga-top/skyline-foundation/quicserver/quic/utils/wire"

type Frame struct {
	wire.Frame
	OnLost  func(wire.Frame)
	OnAcked func(wire.Frame)
}
