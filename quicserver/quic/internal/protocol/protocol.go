package protocol

type ApplicationErrorCode uint64

type ECN uint8

const (
	ECNNon ECN = iota // 00
	ECT1              // 01
	ECT0              // 10
	ECNCE             // 11
)

type StatelessResetToken [16]byte
// A ByteCount in QUIC
type ByteCount uint64


// MaxReceivePacketSize maximum packet size of any QUIC packet, based on
// ethernet's max size, minus the IP and UDP headers. IPv6 has a 40 byte header,
// UDP adds an additional 8 bytes.  This is a total overhead of 48 bytes.
// Ethernet's max packet size is 1500 bytes,  1500 - 48 = 1452.
const MaxReceivePacketSize ByteCount = 1452