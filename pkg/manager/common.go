package manager

import "github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"

func GetClusterCacheKey(cluster *model.ClusterConfig) string {
	return cluster.Name
}
