package registry

import (
	"context"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
)

//Interface to maintain cluster configuration for a sharding manager identity
type ClusterConfiguration interface {
	GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*model.ClusterConfig, error)
	BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*model.ClusterConfig, error)
}

//Interface to maintain identity configuration for a cluster
type IdentityConfiguration interface {
	GetIdentitiesByCluster(ctx context.Context, clusterName string) ([]*model.IdentityConfig, error)
}

type RegistryClient struct {
	RegistryEndpoint string
}

func (c *RegistryClient) GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*model.ClusterConfig, error) {
	return nil, nil
}

func (c *RegistryClient) BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*model.ClusterConfig, error) {
	return nil, nil
}

func (c *RegistryClient) GetIdentitiesByCluster(ctx context.Context, clusterName string) ([]*model.IdentityConfig, error) {
	return nil, nil
}
