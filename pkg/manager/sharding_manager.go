package manager

import (
	"context"
	"fmt"

	admiralV1 "github.com/istio-ecosystem/admiral-api/pkg/client/clientset/versioned/typed/admiral/v1"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/controller"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
)

type shardingManager struct {
	admiralAPIClient admiralV1.AdmiralV1Interface
	registryClient   registry.RegistryConfigInterface
	cache            model.ShardingMangerCache
	shardHandler     controller.ShardInterface
	identity         string
}

func NewShardingManager(
	ctx context.Context,
	shardHandler controller.ShardInterface,
	client model.Clients,
	identity string) (*shardingManager, error) {
	return &shardingManager{
		cache: model.ShardingMangerCache{
			ClusterCache: []registry.ClusterConfig{},
		},
		admiralAPIClient: client.AdmiralClient,
		registryClient:   client.RegistryClient,
		shardHandler:     shardHandler,
		identity:         identity,
	}, nil
}

func (sm *shardingManager) Start(ctx context.Context) error {
	// Bulk sync initial configurations
	err := sm.bulkSync(ctx)
	if err != nil {
		return fmt.Errorf("unable to bulk sync configurations: %v", err)
	}
	// Derive shard configurations from configurations
	config, err := sm.deriveShardConfiguration()
	if err != nil {
		return fmt.Errorf("unable to derive shard configurations: %v", err)
	}
	// Create/Update Shard CRD
	err = sm.pushShardConfiguration(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to push shard configuration: %v", err)
	}
	go sm.startReconciler()
	return nil
}

func (sm *shardingManager) pushShardConfiguration(ctx context.Context, config model.ShardingMangerCache) error {
	_, err := sm.shardHandler.Create(ctx, config.ClusterCache, "identity", "operatorIdentity")
	return err
}

func (sm *shardingManager) bulkSync(ctx context.Context) error {
	var (
		cache []registry.ClusterConfig
		err   error
	)
	cache, err = sm.registryConfigSyncer(ctx)
	if err != nil {
		return err
	}
	sm.cache.ClusterCache = cache
	return nil
}

func (sm *shardingManager) deriveShardConfiguration() (model.ShardingMangerCache, error) {
	return sm.cache, nil
}

func (sm *shardingManager) startReconciler() {
}
