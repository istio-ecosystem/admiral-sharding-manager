package model

import (
	admiralv1 "github.com/istio-ecosystem/admiral-api/pkg/client/clientset/versioned/typed/admiral/v1"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
)

type ShardingManagerParams struct {
	ShardingManagerIdentity string
	OperatorIdentityLabel   string
	ShardIdentityLabel      string
	ShardNamespace          string
	KubeconfigPath          string
	RegistryEndpoint        string
}

type ShardingManagerConfig struct {
	AdmiralApiClient admiralv1.AdmiralV1Interface
	RegistryClient   registry.RegistryConfigInterface
	Cache            ShardingMangerCache
}

type ShardingMangerCache struct {
	ClusterCache []registry.ClusterConfig
}
