package controller

import (
	"context"
	typeV1 "github.com/istio-ecosystem/admiral-api/pkg/apis/admiral/v1"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Interface to manage shards
type ShardInteface interface {
	// create shard resource on a kubernetes cluster
	Create(ctx context.Context, clusterConfiguration interface{}) (*typeV1.Shard, error)
	// update shard resource on a kubernetes cluster
	Update(ctx context.Context, clusterConfiguration interface{}) (*typeV1.Shard, error)
	// delete shard resource on a kubernetes cluster
	Delete(ctx context.Context, shard *typeV1.Shard) error
}

type shardHandler struct {
	config         *model.ShardingManagerConfig
	shardNamespace string
}

// initializes ShardHandler with sharding manager configuration and shard namespace
func NewShardHandler(ShardConfig *model.ShardingManagerConfig, namespace string) shardHandler {
	shardHandler := shardHandler{
		config:         ShardConfig,
		shardNamespace: namespace,
	}

	return shardHandler
}

func (sh *shardHandler) Create(ctx context.Context, clusterConfiguration interface{}) (*typeV1.Shard, error) {
	clusterConfig := clusterConfiguration.(registry.ShardClusterConfig)

	shardToCreate := buildShardResource(clusterConfig.Clusters)

	createdShard, err := sh.config.AdmiralApiClient.Shards(sh.shardNamespace).Create(ctx, shardToCreate, metav1.CreateOptions{})
	if err != nil {
		log.Error("failed to create shard resource")
	}

	return createdShard, err
}

func (sh *shardHandler) Update(ctx context.Context, clusterConfiguration interface{}) (*typeV1.Shard, error) {
	clusterConfig := clusterConfiguration.(registry.ShardClusterConfig)

	shardToUpdate := buildShardResource(clusterConfig.Clusters)

	updatedShard, err := sh.config.AdmiralApiClient.Shards(sh.shardNamespace).Update(ctx, shardToUpdate, metav1.UpdateOptions{})
	if err != nil {
		log.Error("failed to update shard resource")
	}
	return updatedShard, err
}

func (sh *shardHandler) Delete(ctx context.Context, shard *typeV1.Shard) error {
	err := sh.config.AdmiralApiClient.Shards(sh.shardNamespace).Delete(ctx, shard.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Error("failed to delete shard resource")
	}
	return err
}

func buildShardResource(clusterConfigs []registry.ClusterConfig) *typeV1.Shard {
	var clusters []typeV1.ClusterShards

	for _, clusterConfig := range clusterConfigs {
		cluster := typeV1.ClusterShards{
			Name:     clusterConfig.Name,
			Locality: clusterConfig.Locality,
		}

		var identities []typeV1.IdentityItem
		for _, identityConfig := range clusterConfig.IdentityConfig.IdentityMetadata {
			identity := typeV1.IdentityItem{
				Name:        identityConfig.Name,
				Environment: identityConfig.Environment,
			}
			identities = append(identities, identity)
		}
		cluster.Identities = identities
		clusters = append(clusters, cluster)
	}

	shard := &typeV1.Shard{
		ObjectMeta: metav1.ObjectMeta{},
		Spec: typeV1.ShardSpec{
			Clusters: clusters,
		},
	}

	return shard
}

func handleLoadDistribution(ctx context.Context) error {
	return nil
}
