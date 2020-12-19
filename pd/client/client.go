package client

import (
	"context"
	"github.com/pingcap/kvproto/pkg/metapb"
	"github.com/pingcap/kvproto/pkg/pdpb"
	"time"
)

const (
	defaultPDTimeout      = 3 * time.Second
	dialTimeout           = 3 * time.Second
	updateLeaderTimeout   = time.Second // Use a shorter timeout to recover faster from network isolation.
	tsLoopDCCheckInterval = time.Minute
	maxMergeTSORequests   = 10000 // should be higher if client is sending requests in burst
	maxInitClusterRetries = 100
)

type Client interface {

	// GetClusterID gets the cluster ID from PD.
	GetClusterID(ctx context.Context) uint64
	// GetAllMembers gets the members Info from PD
	GetAllMembers(ctx context.Context) ([]*pdpb.Member, error)
	// GetLeaderAddr returns current leader's address. It returns "" before
	// syncing leader from server.
	GetLeaderAddr() string
	// GetTS gets a timestamp from PD.
	GetTS(ctx context.Context) (int64, int64, error)
	// GetTSAsync gets a timestamp from PD, without block the caller.
	GetTSAsync(ctx context.Context) TSFuture
	// GetLocalTS gets a local timestamp from PD.
	GetLocalTS(ctx context.Context, dcLocation string) (int64, int64, error)
	// GetLocalTSAsync gets a local timestamp from PD, without block the caller.
	GetLocalTSAsync(ctx context.Context, dcLocation string) TSFuture
	// GetRegion gets a region and its leader Peer from PD by key.
	// The region may expire after split. Caller is responsible for caching and
	// taking care of region change.
	// Also it may return nil if PD finds no Region for the key temporarily,
	// client should retry later.
	GetRegion(ctx context.Context, key []byte) (*Region, error)
	// GetRegionFromMember gets a region from certain members.
	GetRegionFromMember(ctx context.Context, key []byte, memberURLs []string) (*Region, error)
	// GetPrevRegion gets the previous region and its leader Peer of the region where the key is located.
	GetPrevRegion(ctx context.Context, key []byte) (*Region, error)
	// GetRegionByID gets a region and its leader Peer from PD by id.
	GetRegionByID(ctx context.Context, regionID uint64) (*Region, error)
	// ScanRegion gets a list of regions, starts from the region that contains key.
	// Limit limits the maximum number of regions returned.
	// If a region has no leader, corresponding leader will be placed by a peer
	// with empty value (PeerID is 0).
	ScanRegions(ctx context.Context, key, endKey []byte, limit int) ([]*Region, error)
	// GetStore gets a store from PD by store id.
	// The store may expire later. Caller is responsible for caching and taking care
	// of store change.
	GetStore(ctx context.Context, storeID uint64) (*metapb.Store, error)
	// GetAllStores gets all stores from pd.
	// The store may expire later. Caller is responsible for caching and taking care
	// of store change.
	GetAllStores(ctx context.Context, opts ...GetStoreOption) ([]*metapb.Store, error)
	// Update GC safe point. TiKV will check it and do GC themselves if necessary.
	// If the given safePoint is less than the current one, it will not be updated.
	// Returns the new safePoint after updating.
	UpdateGCSafePoint(ctx context.Context, safePoint uint64) (uint64, error)

	// UpdateServiceGCSafePoint updates the safepoint for specific service and
	// returns the minimum safepoint across all services, this value is used to
	// determine the safepoint for multiple services, it does not trigger a GC
	// job. Use UpdateGCSafePoint to trigger the GC job if needed.
	UpdateServiceGCSafePoint(ctx context.Context, serviceID string, ttl int64, safePoint uint64) (uint64, error)
	// ScatterRegion scatters the specified region. Should use it for a batch of regions,
	// and the distribution of these regions will be dispersed.
	// NOTICE: This method is the old version of ScatterRegions, you should use the later one as your first choice.
	ScatterRegion(ctx context.Context, regionID uint64) error
	// ScatterRegions scatters the specified regions. Should use it for a batch of regions,
	// and the distribution of these regions will be dispersed.
	ScatterRegions(ctx context.Context, regionsID []uint64, opts ...RegionsOption) (*pdpb.ScatterRegionResponse, error)
	// SplitRegions split regions by given split keys
	SplitRegions(ctx context.Context, splitKeys [][]byte, opts ...RegionsOption) (*pdpb.SplitRegionsResponse, error)
	// GetOperator gets the status of operator of the specified region.
	GetOperator(ctx context.Context, regionID uint64) (*pdpb.GetOperatorResponse, error)
	// Close closes the client.
	Close()
}