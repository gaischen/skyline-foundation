package wire

import "time"

type AckFrame struct {
	AckRanges []AckRange // has to be ordered. The highest ACK range goes first, the lowest ACK range goes last
	DelayTime time.Duration

	ECT0, ECT1, ECNCE uint64
}

