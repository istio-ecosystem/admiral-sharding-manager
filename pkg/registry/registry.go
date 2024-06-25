package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// interface to interact with registry service to maintain resource configuration
type RegistryConfigInterface interface {
	// fetch cluster configuration by sharding manager identity
	GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentity string) (interface{}, error)
	// bulk fetch cluster configuration by sharding manager identity
	BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentity string) (interface{}, error)
	// fetch identities by cluster name
	GetIdentitiesByCluster(ctx context.Context, clusterName string) (interface{}, error)
}

type registryClient struct {
	registryEndpoint string
}

type ShardClusterConfig struct {
	Clusters        []clusterConfig `json:"clusters, omitempty"`
	LastUpdatedTime string          `json:"lastUpdatedTime, omitempty"`
	ResourceVersion string          `json:"resourceVersion, omitempty"`
}

// cluster configuration for sharding manager identity
type clusterConfig struct {
	Name           string          `json:"name,omitempty"`
	Locality       string          `json:"locality,omitempty"`
	Metadata       clusterMetadata `json:"metadata,omitempty"`
	IdentityConfig identityConfig  `json:"assets,omitempty"`
}

type clusterMetadata struct {
}

// mesh workload identity configuration for cluster
type identityConfig struct {
	ClusterName      string             `json:"clustername,omitempty"`
	IdentityMetadata []identityMetadata `json:"assetMetadata,omitempty"`
}

type identityMetadata struct {
	Name             string `json:"asset,omitempty"`
	SourceAsset      bool   `json:"sourceAsset,omitempty"`
	DestinationAsset bool   `json:"destinationAsset,omitempty"`
}

// initializes registry client configuration
func NewRegistryClient(endpoint string) RegistryConfigInterface {
	return &registryClient{
		registryEndpoint: endpoint,
	}
}

func (c *registryClient) GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentity string) (interface{}, error) {
	var clusterConfigData ShardClusterConfig
	tid := uuid.NewString()

	ctxLogger := log.WithFields(log.Fields{
		"smIdentity": shardingManagerIdentity,
		"tid":        tid,
	})
	ctxLogger.Infof("Get cluster configuration for provided sharding manager identity")

	err := checkIfRegistryClientIsInitailized(c)
	if err != nil {
		ctxLogger.WithError(err).Error("registry client not initialized")
		return clusterConfigData, err
	}

	byteValue, err := os.ReadFile("../testdata/" + strings.TrimSpace(shardingManagerIdentity) + ".json")
	if err != nil {
		ctxLogger.WithError(err).Error("failed to get cluster configuration from registry")
		return clusterConfigData, err
	}

	err = json.Unmarshal(byteValue, &clusterConfigData)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to unmarshal cluster configuration")
		return clusterConfigData, err
	}

	return clusterConfigData, nil
}

func checkIfRegistryClientIsInitailized(registryClient *registryClient) error {
	if registryClient == nil || registryClient.registryEndpoint == "" {
		return fmt.Errorf("registry client is not initialized")
	}
	return nil
}

func (c *registryClient) BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentity string) (interface{}, error) {
	var clusterConfigData ShardClusterConfig
	tid := uuid.NewString()

	ctxLogger := log.WithFields(log.Fields{
		"smIdentity": shardingManagerIdentity,
		"tid":        tid,
	})
	ctxLogger.Infof("bulk sync cluster configuration for provided sharding manager identity")

	err := checkIfRegistryClientIsInitailized(c)
	if err != nil {
		ctxLogger.WithError(err).Error("registry client not initialized")
		return clusterConfigData, err
	}

	byteValue, err := os.ReadFile("../testdata/" + shardingManagerIdentity + "-bulk.json")
	if err != nil {
		ctxLogger.WithError(err).Error("failed perform bulk sync for cluster configuration from registry")
		return clusterConfigData, err
	}

	err = json.Unmarshal(byteValue, &clusterConfigData)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to unmarshal cluster configuration")
		return clusterConfigData, err
	}

	return clusterConfigData, nil
}

func (c *registryClient) GetIdentitiesByCluster(ctx context.Context, clusterName string) (interface{}, error) {
	var identityConfig identityConfig
	tid := uuid.NewString()

	ctxLogger := log.WithFields(log.Fields{
		"clusterName": clusterName,
		"tid":         tid,
	})
	ctxLogger.Infof("Get cluster configuration for provided sharding manager identity")

	err := checkIfRegistryClientIsInitailized(c)
	if err != nil {
		ctxLogger.WithError(err).Error("registry client not initialized")
		return identityConfig, err
	}

	byteValue, err := os.ReadFile("../testdata/" + clusterName + ".json")
	if err != nil {
		ctxLogger.WithError(err).Error("failed to get cluster configuration from registry")
		return identityConfig, err
	}

	err = json.Unmarshal(byteValue, &identityConfig)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to unmarshal cluster configuration")
		return identityConfig, err
	}

	return identityConfig, nil
}
