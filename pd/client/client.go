package client

import "time"

const (
	defaultPDTimeout      = 3 * time.Second
	dialTimeout           = 3 * time.Second
	updateLeaderTimeout   = time.Second // Use a shorter timeout to recover faster from network isolation.
	tsLoopDCCheckInterval = time.Minute
	maxMergeTSORequests   = 10000 // should be higher if client is sending requests in burst
	maxInitClusterRetries = 100
)
