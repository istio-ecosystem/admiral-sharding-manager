package manager

import (
	"context"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
	"github.com/sirupsen/logrus"
)

func InitializeShardingManager(ctx context.Context, params *model.ShardingManagerParams) (*model.ShardingManagerConfig, error) {

	smConfig := &model.ShardingManagerConfig{}
	var err error

	//setup admiral client
	kubeClient := KubeClient{}
	smConfig.AdmiralApiClient, err = kubeClient.LoadAdmiralApiClientFromPath(params.KubeconfigPath)
	if err != nil {
		logrus.Error("failed to initialize admiral api client")
	}
	//setup registry client
	smConfig.RegistryClient = registry.RegistryClient{}

	//TODO: setup oms client and subscribe to topic specific for this sharding manager identity

	smConfig.Cache = model.ShardingMangerCache{}

	err = LoadRegistryConfiguration(ctx, smConfig, params)
	if err != nil {
		logrus.Error("failed to initialize registry client")
	}

	return smConfig, err
}

func LoadRegistryConfiguration(ctx context.Context, config *model.ShardingManagerConfig, params *model.ShardingManagerParams) error {
	var err error

	shardClusters, err := config.RegistryClient.GetClustersByShardingManagerIdentity(ctx, params.ShardingManagerIdentity)
	if err != nil {
		return err
	}

	for _, cluster := range shardClusters {
		key := GetClusterCacheKey(cluster)
		clusterIdentities, err := config.RegistryClient.GetIdentitiesByCluster(ctx, cluster.Name)
		if err != nil {
			return err
		}

		config.Cache.IdentityCache.Store(key, clusterIdentities)
	}
	return nil
}
