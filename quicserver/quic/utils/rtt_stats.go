package utils

import "time"

// RTTStats provides round-trip statistics
type RTTStats struct {
	hasMeasurement bool

	minRTT        time.Duration
	latestRTT     time.Duration
	smoothedRTT   time.Duration
	meanDeviation time.Duration

	maxAckDelay time.Duration
}
