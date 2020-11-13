package keepalive

import "time"

type ServerParameters struct {
	MaxConnectionIdle     time.Duration
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
	Time                  time.Duration
	Timeout               time.Duration
}

type EnforcementPolicy struct {
	MinTime time.Duration
	//if true server希望在active stream==0的时候还想用ping保持连接
	PermitWithoutStream bool
}
