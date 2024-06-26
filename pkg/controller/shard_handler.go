package controller

import (
	"context"
	"fmt"
	typeV1 "github.com/istio-ecosystem/admiral-api/pkg/apis/admiral/v1"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SHARD_IDENTITY_LABEL = "admiral.io/shardIdentity"
)

// Interface to manage shards
type ShardInteface interface {
	// create shard resource on a kubernetes cluster
	Create(ctx context.Context, clusterConfiguration []registry.ClusterConfig, shardName string, operatorIdentity string) (*typeV1.Shard, error)
	// update shard resource on a kubernetes cluster
	Update(ctx context.Context, clusterConfiguration []registry.ClusterConfig, shardName string, operatorIdentity string) (*typeV1.Shard, error)
	// delete shard resource on a kubernetes cluster
	Delete(ctx context.Context, shard *typeV1.Shard) error
}

type shardHandler struct {
	config *model.ShardingManagerConfig
	params *model.ShardingManagerParams
}

// initializes ShardHandler with sharding manager configuration and shard namespace
func NewShardHandler(ShardConfig *model.ShardingManagerConfig, smParams *model.ShardingManagerParams) shardHandler {
	shardHandler := shardHandler{
		config: ShardConfig,
		params: smParams,
	}

	return shardHandler
}

func (sh *shardHandler) Create(ctx context.Context, clusterConfiguration []registry.ClusterConfig, shardName string, operatorIdentity string) (*typeV1.Shard, error) {
	if sh.config.AdmiralApiClient == nil {
		return nil, fmt.Errorf("admiral api client is not initialized")
	}

	shardToCreate := buildShardResource(clusterConfiguration, sh.params, shardName, operatorIdentity)

	createdShard, err := sh.config.AdmiralApiClient.Shards(sh.params.ShardNamespace).Create(ctx, shardToCreate, metav1.CreateOptions{})
	if err != nil {
		log.Error("failed to create shard resource")
	}
	return createdShard, err
}

func (sh *shardHandler) Update(ctx context.Context, clusterConfiguration []registry.ClusterConfig, shardName string, operatorIdentity string) (*typeV1.Shard, error) {
	var updatedShard *typeV1.Shard

	if sh.config.AdmiralApiClient == nil {
		return nil, fmt.Errorf("admiral api client is not initialized")
	}

	existingShard, err := sh.config.AdmiralApiClient.Shards(sh.params.ShardNamespace).Get(ctx, shardName, metav1.GetOptions{})
	shardToUpdate := buildShardResource(clusterConfiguration, sh.params, shardName, operatorIdentity)

	if existingShard != nil && shardToUpdate != nil {
		existingShard.Labels = updatedShard.Labels
		existingShard.Annotations = updatedShard.Annotations
		existingShard.Spec = updatedShard.Spec

		updatedShard, err = sh.config.AdmiralApiClient.Shards(sh.params.ShardNamespace).Update(ctx, shardToUpdate, metav1.UpdateOptions{})
		if err != nil {
			log.Error("failed to update shard resource")
		}
	}
	return updatedShard, err
}

func (sh *shardHandler) Delete(ctx context.Context, shard *typeV1.Shard) error {
	if sh.config.AdmiralApiClient == nil {
		return fmt.Errorf("admiral api client is not initialized")
	}

	err := sh.config.AdmiralApiClient.Shards(sh.params.ShardNamespace).Delete(ctx, shard.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Error("failed to delete shard resource")
	}
	return err
}

func buildShardResource(clusterConfigs []registry.ClusterConfig, smParam *model.ShardingManagerParams, shardName string, operatorIdentity string) *typeV1.Shard {
	var clusters []typeV1.ClusterShards

	labels := make(map[string]string)
	labels[SHARD_IDENTITY_LABEL] = smParam.ShardingManagerIdentity
	labels[smParam.OperatorIdentityLabel] = operatorIdentity

	for _, clusterConfig := range clusterConfigs {
		cluster := typeV1.ClusterShards{
			Name:     clusterConfig.Name,
			Locality: clusterConfig.Locality,
		}

		var identities []typeV1.IdentityItem
		for _, identityConfig := range clusterConfig.IdentityConfig.AssetList {
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      shardName,
			Namespace: smParam.ShardNamespace,
			Labels:    labels,
		},
		Spec: typeV1.ShardSpec{
			Clusters: clusters,
		},
	}

	return shard
}

// distributed clusterconfig into shard resource
// TODO: Currently does not have logic to distribute cluster configuration, in next phase will have the load distribution in place
func (sh *shardHandler) HandleLoadDistribution(ctx context.Context) error {
	operatorIdentity := "0-1"
	shardName := "shard-" + operatorIdentity

	_, err := sh.Create(ctx, sh.config.Cache.ClusterCache, shardName, operatorIdentity)

	if err != nil {
		return err
	}

	return nil
}
