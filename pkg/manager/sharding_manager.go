package manager

import (
	"context"
	"fmt"

	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
	log "github.com/sirupsen/logrus"
)

// initializes sharding manager with required kubernetes clients, registry client and bootstrap configuration
func InitializeShardingManager(ctx context.Context, params *model.ShardingManagerParams) (*model.ShardingManagerConfig, error) {

	smConfig := &model.ShardingManagerConfig{}
	var err error

	//setup admiral client
	var kubeClient LoadKubeClient = &kubeClient{}

	smConfig.AdmiralApiClient, err = kubeClient.LoadAdmiralApiClientFromPath(params.KubeconfigPath)
	if err != nil {
		log.Error("failed to initialize admiral api client")
	}

	//setup registry client
	//TODO: send registry endpoint
	smConfig.RegistryClient = registry.NewRegistryClient(registry.WithEndpoint(params.RegistryEndpoint))

	//TODO: setup oms client and subscribe to topic specific for this sharding manager identity

	smConfig.Cache = model.ShardingMangerCache{
		ClusterCache: []registry.ClusterConfig{},
	}
	err = registryConfigSyncer(ctx, smConfig, params)
	if err != nil {
		log.Error("failed to initialize registry client")
	}

	return smConfig, err
}

// loads configuration from registry for provide sharding manager identity
func registryConfigSyncer(ctx context.Context, config *model.ShardingManagerConfig, params *model.ShardingManagerParams) error {
	var err error
	clusterConfiguration, err := config.RegistryClient.GetClustersByShardingManagerIdentity(ctx, params.ShardingManagerIdentity)
	if err != nil {
		return err
	}

	if clusterConfiguration == nil {
		return fmt.Errorf("failed to get cluster configuration from registry")
	}

	for _, cluster := range clusterConfiguration.(registry.ShardClusterConfig).Clusters {
		identityConfig, err := config.RegistryClient.GetIdentitiesByCluster(ctx, cluster.Name)

		if err != nil {
			return err
		}

		cluster.IdentityConfig = identityConfig.(registry.IdentityConfig)
		config.Cache.ClusterCache = append(config.Cache.ClusterCache, cluster)
	}

	return nil
}
