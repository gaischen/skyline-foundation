package quic

type sendQueue struct {
	queue       chan *packetBuffer
	closeCalled chan struct{}
	runStopped  chan struct{}
	conn        sendConn
}


