package quic

type packetBuffer struct {
	Data []byte
	// refCount counts how many packets Data is used in.
	// It doesn't support concurrent use.
	// It is > 1 when used for coalesced packet.
	refCount int
}
