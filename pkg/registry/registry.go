package registry

import (
	"context"
)

// Interface to maintain registry configuration for a sharding manager identity
type RegistryConfigInterface interface {
	GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*clusterConfig, error)
	BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*clusterConfig, error)
	GetIdentitiesByCluster(ctx context.Context, clusterName string) ([]*identityConfig, error)
}

type registryClient struct {
	registryEndpoint string
}

// cluster configuration for sharding manager identity
type clusterConfig struct {
	Name     string          `json:"name,omitempty"`
	Locality string          `json:"locality,omitempty"`
	Metadata clusterMetadata `json:"metadata,omitempty"`
}

type clusterMetadata struct {
}

// mesh workload identity configuration for cluster
type identityConfig struct {
	ClusterName      string           `json:"clustername,omitempty"`
	IdentityMetadata identityMetadata `json:"assetMetadata,omitempty"`
}

type identityMetadata struct {
	Name             string `json:"asset,omitempty"`
	SrouceAsset      string `json:"sourceAsset,omitempty"`
	DestinationAsset string `json:"destinationAsset,omitempty"`
}

func NewRegistryClient(endpoint string) RegistryConfigInterface {
	return &registryClient{
		registryEndpoint: endpoint,
	}
}

func (c *registryClient) GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*clusterConfig, error) {
	return nil, nil
}

func (c *registryClient) BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*clusterConfig, error) {
	return nil, nil
}

func (c *registryClient) GetIdentitiesByCluster(ctx context.Context, clusterName string) ([]*identityConfig, error) {
	return nil, nil
}
