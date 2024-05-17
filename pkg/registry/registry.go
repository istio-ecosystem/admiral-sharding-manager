package registry

import (
	"context"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
)

//Interface to maintain cluster configuration for a sharding manager identity
type ClusterConfiguration interface {
	GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]model.Cluster, error)
	BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]model.Cluster, error)
}

//Interface to maintain identity configuration for a cluster
type IdentityConfiguration interface {
	GetIdentitiesByCluster(ctx context.Context, clusterName string) ([]model.Identity, error)
}
