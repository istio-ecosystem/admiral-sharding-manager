package manager

import (
	"context"
	"fmt"

	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
)

// loads configuration from registry for provide sharding manager identity
func (sm *shardingManager) registryConfigSyncer(ctx context.Context) ([]registry.ClusterConfig, error) {
	var err error
	var cache []registry.ClusterConfig
	clusterConfiguration, err := sm.registryClient.GetClustersByShardingManagerIdentity(ctx, sm.identity)
	if err != nil {
		return cache, err
	}
	for _, cluster := range clusterConfiguration.Clusters {
		identityConfig, err := sm.registryClient.GetIdentitiesByCluster(ctx, cluster.Name)
		if err != nil {
			return cache, err
		}
		cluster.IdentityConfig = identityConfig
		cache = append(cache, cluster)
	}
	fmt.Printf("cache: %+v", cache)
	return cache, nil
}
