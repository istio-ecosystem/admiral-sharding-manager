package registry

import (
	"context"
)

// interface to interact with registry service to maintain resource configuration
type RegistryConfigInterface interface {
	// fetch cluster configuration by sharding manager identity
	GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*clusterConfig, error)
	// bulk fetch cluster configuration by sharding manager identity
	BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentityName string) ([]*clusterConfig, error)
	// fetch identities by cluster name
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
	SourceAsset      string `json:"sourceAsset,omitempty"`
	DestinationAsset string `json:"destinationAsset,omitempty"`
}

// initializes registry client configuration
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
