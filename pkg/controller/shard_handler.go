package controller

import (
	"context"
	typeV1 "github.com/istio-ecosystem/admiral-api/pkg/apis/admiral/v1"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ShardInteface interface {
	CreateOrUpdateShard(ctx context.Context, shard *typeV1.Shard) (*typeV1.Shard, error)
	DeleteShard(ctx context.Context, shard *typeV1.Shard) error
}

type ShardHandler struct {
	Config         *model.ShardingManagerConfig
	ShardNamespace string
}

func NewShardHandler(ShardConfig *model.ShardingManagerConfig, namespace string) ShardHandler {
	shardHandler := ShardHandler{
		Config:         ShardConfig,
		ShardNamespace: namespace,
	}

	return shardHandler
}

func (sh ShardHandler) CreateOrUpdateShard(ctx context.Context, shard *typeV1.Shard) (*typeV1.Shard, error) {
	shard = buildShardResource(ctx)
	updatedShard, err := sh.Config.AdmiralApiClient.Shards(sh.ShardNamespace).Create(ctx, shard, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("failed to create or update shard resource")
	}
	return updatedShard, err
}

func (sh ShardHandler) DeleteShard(ctx context.Context, shard *typeV1.Shard) error {
	err := sh.Config.AdmiralApiClient.Shards(sh.ShardNamespace).Delete(ctx, shard.Name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("failed to delete shard resource")
	}
	return err
}

func buildShardResource(ctx context.Context) *typeV1.Shard {
	//TODO: Build the shard resource from the cache

	var identities []typeV1.IdentityItem
	identity := typeV1.IdentityItem{
		Name:        "",
		Environment: "",
	}
	identities = append(identities, identity)

	var clusters []typeV1.ClusterShards

	cluster := typeV1.ClusterShards{
		Name:       "",
		Locality:   "",
		Identities: identities,
	}

	clusters = append(clusters, cluster)

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
