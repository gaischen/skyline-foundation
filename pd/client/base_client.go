package client

import (
	"context"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
	"time"
)

type baseClient struct {
	urls        []string
	clusterID   uint64
	leader      atomic.Value
	clientConns sync.Map
	allocators  sync.Map

	checkLeaderCh chan struct{}

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	security        SecurityOption
	gRPCDialOptions []grpc.DialOption
	timeout         time.Duration
	maxRetryTimes   int
}



type SecurityOption struct {
	CAPath   string
	CertPath string
	KeyPath  string
}
