package controller

import (
	"context"
	"fmt"

	typeV1 "github.com/istio-ecosystem/admiral-api/pkg/apis/admiral/v1"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ShardIdentity = "admiral.io/shardIdentity"
)

// Interface to manage shards
type ShardInterface interface {
	// create shard resource on a kubernetes cluster
	Create(ctx context.Context, clusterConfiguration []registry.ClusterConfig, shardName string, operatorIdentity string) (*typeV1.Shard, error)
	// update shard resource on a kubernetes cluster
	Update(ctx context.Context, clusterConfiguration []registry.ClusterConfig, shardName string, operatorIdentity string) (*typeV1.Shard, error)
	// delete shard resource on a kubernetes cluster
	Delete(ctx context.Context, shard *typeV1.Shard) error
}

type shardHandler struct {
	clients model.Clients
	params  *model.ShardingManagerParams
}

// initializes ShardHandler with sharding manager configuration and shard namespace
func NewShardHandler(clients model.Clients, smParams *model.ShardingManagerParams) *shardHandler {
	shardHandler := &shardHandler{
		clients: clients,
		params:  smParams,
	}
	return shardHandler
}

func (sh *shardHandler) Create(
	ctx context.Context,
	clusterConfiguration []registry.ClusterConfig,
	shardName string,
	operatorIdentity string) (*typeV1.Shard, error) {
	shardToCreate := buildShardResource(clusterConfiguration, sh.params, shardName, operatorIdentity)
	return sh.clients.AdmiralClient.Shards(sh.params.ShardNamespace).Create(ctx, shardToCreate, metav1.CreateOptions{})
}

func (sh *shardHandler) Update(
	ctx context.Context,
	clusterConfiguration []registry.ClusterConfig,
	shardName string,
	operatorIdentity string) (*typeV1.Shard, error) {
	var updatedShard *typeV1.Shard
	existingShard, err := sh.clients.AdmiralClient.Shards(sh.params.ShardNamespace).Get(ctx, shardName, metav1.GetOptions{})
	shardToUpdate := buildShardResource(clusterConfiguration, sh.params, shardName, operatorIdentity)

	if existingShard != nil && shardToUpdate != nil {
		existingShard.Labels = shardToUpdate.Labels
		existingShard.Annotations = shardToUpdate.Annotations
		existingShard.Spec = shardToUpdate.Spec

		updatedShard, err = sh.clients.AdmiralClient.Shards(sh.params.ShardNamespace).Update(ctx, existingShard, metav1.UpdateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to update shard resource: %v", err)
		}
	}
	return updatedShard, err
}

func (sh *shardHandler) Delete(ctx context.Context, shard *typeV1.Shard) error {
	err := sh.clients.AdmiralClient.Shards(sh.params.ShardNamespace).Delete(ctx, shard.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete shard resource: %v", err)
	}
	return err
}

func buildShardResource(clusterConfigs []registry.ClusterConfig, smParam *model.ShardingManagerParams, shardName string, operatorIdentity string) *typeV1.Shard {
	var (
		clusters []typeV1.ClusterShards
		labels   = make(map[string]string)
	)
	labels[ShardIdentity] = smParam.ShardingManagerIdentity
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
