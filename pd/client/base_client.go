package client

import "C"
import (
	"context"
	"github.com/pingcap/kvproto/pkg/pdpb"
	"github.com/pkg/errors"
	"github.com/vanga-top/skyline-foundation/log"
	"github.com/vanga-top/skyline-foundation/log/level"
	"github.com/vanga-top/skyline-foundation/pd/pkg/grpcutil"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

var logger = log.NewLogger("pd", level.ERROR)

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

type ClientOption func(c *baseClient)

type SecurityOption struct {
	CAPath   string
	CertPath string
	KeyPath  string
}

func newBaseClient(ctx context.Context, urls []string, security SecurityOption, opts ...ClientOption) (*baseClient, error) {
	ctx1, cancel := context.WithCancel(ctx)
	c := &baseClient{
		urls:          urls,
		checkLeaderCh: make(chan struct{}, 1),
		ctx:           ctx1,
		cancel:        cancel,
		security:      security,
		timeout:       defaultPDTimeout,
		maxRetryTimes: maxInitClusterRetries,
	}
	for _, opt := range opts {
		opt(c)
	}
	//init retry
	if err := c.initRetry(c.initClusterID); err != nil {
		c.cancel()
		return nil, err
	}

	if err := c.initRetry(c.updateLeader); err != nil {
		c.cancel()
		return nil, err
	}
	c.wg.Add(1)
	go c.leaderLoop()
	return c, nil
}

func (c *baseClient) leaderLoop() {
	defer c.wg.Done()
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	for {
		select {
		case <-c.checkLeaderCh:
		case <-time.After(time.Minute):
		case <-ctx.Done():
			return
		}

		if err := c.updateLeader(); err != nil {
			logger.Error("[pd] failed update leader...", err)
		}
	}
}

func (c *baseClient) updateLeader() error {
	for _, u := range c.urls {
		ctx, cancel := context.WithTimeout(c.ctx, updateLeaderTimeout)
		members, err := c.getMembers(ctx, u)
		if err != nil {
			logger.Warn("[PD] CANNOT UPDATE LEADER..")
		}
		cancel()
		if err := c.switchTSOAllocatorLeader(members.GetTsoAllocatorLeaders()); err != nil {
			return err
		}
		if members.GetLeader() == nil || len(members.GetLeader().GetClientUrls()) == 0 {
			select {
			case <-c.ctx.Done():
				return errors.New("timeout")
			default:
				continue
			}
		}
		c.updateURLs(members.GetMembers())
		return c.switchLeader(members.GetLeader().GetClientUrls())
	}
	return errors.New("updateLeader error...")
}

func (c *baseClient) initRetry(f func() error) error {
	var err error
	for i := 0; i < c.maxRetryTimes; i++ {
		if err = f(); err == nil {
			return nil
		}
		select {
		case <-c.ctx.Done():
			return err
		case <-time.After(time.Second):
		}
	}
	return errors.WithStack(err)
}

func (c *baseClient) initClusterID() error {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	for _, u := range c.urls {
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, c.timeout)
		members, err := c.getMembers(timeoutCtx, u)
		timeoutCancel()
		if err != nil || members.GetHeader() == nil {
			logger.Warn("[pd] failed to get cluster id...")
			continue
		}
		c.clusterID = members.GetHeader().GetClusterId()
		return nil
	}
	return errors.New("failed to get cluster id..")
}

func (c *baseClient) switchTSOAllocatorLeader(allocatorMap map[string]*pdpb.Member) error {
	if len(allocatorMap) == 0 {
		return nil
	}
	for dcLocation, member := range allocatorMap {
		if len(member.GetClientUrls()) == 0 {
			continue
		}
		addr := member.GetClientUrls()[0]
		oldAddr, exist := c.getAllocatorLeaderAddrByDBLocation(dcLocation)
		if exist && addr == oldAddr {
			continue
		}
		logger.Info("[pd] switch dc tso allocator leader",
			zap.String("dc-location", dcLocation),
			zap.String("new-leader", addr),
			zap.String("old-leader", oldAddr))
		if _, err := c.getOrCreateGRPCConn(addr); err != nil {
			return err
		}
		c.allocators.Store(dcLocation, addr)
	}
	c.gcAllocatorLeaderAddr(allocatorMap)
	return nil
}

func (c *baseClient) getAllocatorLeaderAddrByDBLocation(dcLocation string) (string, bool) {
	url, exist := c.allocators.Load(dcLocation)
	if !exist {
		return "", false
	}
	return url.(string), true
}

func (c *baseClient) getMembers(ctx context.Context, url string) (*pdpb.GetMembersResponse, error) {
	cc, err := c.getOrCreateGRPCConn(url)
	if err != nil {
		return nil, err
	}
	member, err := pdpb.NewPDClient(cc).GetMembers(ctx, &pdpb.GetMembersRequest{})
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (c *baseClient) getOrCreateGRPCConn(addr string) (*grpc.ClientConn, error) {
	conn, ok := c.clientConns.Load(addr)
	if ok {
		return conn.(*grpc.ClientConn), nil
	}
	tlsCfg, err := grpcutil.TLSConfig{
		CAPath:   c.security.CAPath,
		CertPath: c.security.CertPath,
		KeyPath:  c.security.KeyPath,
	}.ToTLSConfig()
	if err != nil {
		return nil, err
	}
	dCtx, cancel := context.WithTimeout(c.ctx, dialTimeout)
	defer cancel()
	cc, err := grpcutil.GetClientConn(dCtx, addr, tlsCfg, c.gRPCDialOptions...)
	if err != nil {
		return nil, err
	}
	if old, ok := c.clientConns.Load(addr); ok {
		cc.Close()
		return old.(*grpc.ClientConn), nil
	}
	c.clientConns.Store(addr, cc)
	return cc, nil
}

func (c *baseClient) gcAllocatorLeaderAddr(curAllocatorMap map[string]*pdpb.Member) {
	c.allocators.Range(func(dcLocation, _ interface{}) bool {
		if dcLocation.(string) == "global" {
			return true
		}
		if _, exist := curAllocatorMap[dcLocation.(string)]; !exist {
			c.allocators.Delete(dcLocation)
		}
		return true
	})
}

func (c *baseClient) updateURLs(members []*pdpb.Member) {
	urls := make([]string, 0, len(members))
	for _, m := range members {
		urls = append(urls, m.GetClientUrls()...)
	}
	sort.Strings(urls)
	if reflect.DeepEqual(c.urls, urls) {
		return
	}
	c.urls = urls
}

func (c *baseClient) switchLeader(addrs []string) error {
	addr := addrs[0]
	oldLeader := c.GetLeaderAddr()
	if addr == oldLeader {
		return nil
	}
	if _, err := c.getOrCreateGRPCConn(addr); err != nil {
		return err
	}
	c.leader.Store(addr)
	c.allocators.Store("global", addr)
	return nil
}

func (c *baseClient) GetLeaderAddr() string {
	leaderAddr := c.leader.Load()
	if leaderAddr == nil {
		return ""
	}
	return leaderAddr.(string)
}
