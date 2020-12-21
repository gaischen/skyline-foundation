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
