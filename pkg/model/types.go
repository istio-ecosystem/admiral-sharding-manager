package model

import (
	admiralv1 "github.com/istio-ecosystem/admiral-api/pkg/client/clientset/versioned/typed/admiral/v1"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
	"sync"
)

//cluster configuration for sharding manager identity
type ClusterConfig struct {
	Name     string          `json:"name,omitempty"`
	Locality string          `json:"locality,omitempty"`
	Metadata ClusterMetadata `json:"metadata,omitempty"`
}

type ClusterMetadata struct {
}

//mesh workload identity configuration for cluster
type IdentityConfig struct {
	ClusterName      string `json:"clustername,omitempty"`
	IdentityMetadata string `json:"assetMetadata,omitempty"`
}

type IdentityMetadata struct {
	Name             string `json:"asset,omitempty"`
	SrouceAsset      string `json:"sourceAsset,omitempty"`
	DestinationAsset string `json:"destinationAsset,omitempty"`
}

type ShardingManagerParams struct {
	ShardingManagerIdentity string
	OperatorIdentityLabel   string
	ShardIdentityLabel      string
	ShardNamespace          string
	KubeconfigPath          string
}

type ShardingManagerConfig struct {
	AdmiralApiClient admiralv1.AdmiralV1Interface
	RegistryClient   registry.RegistryClient
	Cache            ShardingMangerCache
}

type ShardingMangerCache struct {
	IdentityCache *sync.Map
}
