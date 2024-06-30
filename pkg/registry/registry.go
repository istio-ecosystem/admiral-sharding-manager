package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// interface to interact with registry service to maintain resource configuration
type RegistryConfigInterface interface {
	// fetch cluster configuration by sharding manager identity
	GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentity string) (ShardClusterConfig, error)
	// bulk fetch cluster configuration by sharding manager identity
	BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentity string) (ShardClusterConfig, error)
	// fetch identities by cluster name
	GetIdentitiesByCluster(ctx context.Context, clusterName string) (IdentityConfig, error)
}

type registryClient struct {
	registryEndpoint string
}

type ShardClusterConfig struct {
	Clusters        []ClusterConfig `json:"clusters, omitempty"`
	LastUpdatedTime string          `json:"lastUpdatedTime, omitempty"`
	ResourceVersion string          `json:"resourceVersion, omitempty"`
}

// cluster configuration for sharding manager identity
type ClusterConfig struct {
	Name           string          `json:"name,omitempty"`
	Locality       string          `json:"locality,omitempty"`
	Metadata       clusterMetadata `json:"metadata,omitempty"`
	IdentityConfig IdentityConfig  `json:"assets,omitempty"`
}

type clusterMetadata struct {
}

// mesh workload identity configuration for cluster
type IdentityConfig struct {
	ClusterName string      `json:"clustername,omitempty"`
	AssetList   []AssetList `json:"assetList,omitempty"`
}

type AssetList struct {
	Name             string `json:"asset,omitempty"`
	Environment      string `json:"environment,omitempty"`
	SourceAsset      bool   `json:"sourceAsset,omitempty"`
	DestinationAsset bool   `json:"destinationAsset,omitempty"`
}

// initializes registry client configuration
func NewRegistryClient(options ...func(client *registryClient)) *registryClient {
	client := &registryClient{}
	for _, option := range options {
		option(client)
	}
	return client
}

func WithEndpoint(endpoint string) func(client *registryClient) {
	return func(client *registryClient) {
		client.registryEndpoint = endpoint
	}
}

func (c *registryClient) GetClustersByShardingManagerIdentity(ctx context.Context, shardingManagerIdentity string) (ShardClusterConfig, error) {
	var (
		clusterConfigData ShardClusterConfig
		tid               = uuid.NewString()
		ctxLogger         = log.WithFields(log.Fields{
			"smIdentity": shardingManagerIdentity,
			"tid":        tid,
		})
	)
	ctxLogger.Infof("Get cluster configuration for provided sharding manager identity")
	data, err := c.getClustersByShardingManagerIdentity(ctxLogger, shardingManagerIdentity)
	if err != nil {
		return clusterConfigData, fmt.Errorf("unable to fetch config: %v", err)
	}
	err = json.Unmarshal(data, &clusterConfigData)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to unmarshal cluster configuration")
		return clusterConfigData, err
	}
	return clusterConfigData, nil
}

func (c *registryClient) getClustersByShardingManagerIdentity(ctxLogger *logrus.Entry, shardingManagerIdentity string) ([]byte, error) {
	_, base, _, _ := runtime.Caller(0)
	filename := fmt.Sprintf("clusters-for-%s-identity.json", shardingManagerIdentity)
	absPath := filepath.Join(filepath.Dir(base), "/testdata/"+filename)
	data, err := os.ReadFile(absPath)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to get cluster configuration from registry")
		return data, err
	}
	return data, nil
}

func (c *registryClient) BulkSyncByShardingManagerIdentity(ctx context.Context, shardingManagerIdentity string) (ShardClusterConfig, error) {
	var (
		clusterConfigData ShardClusterConfig
		tid               = uuid.NewString()
		ctxLogger         = log.WithFields(log.Fields{
			"smIdentity": shardingManagerIdentity,
			"tid":        tid,
		})
	)
	ctxLogger.Infof("bulk sync cluster configuration for provided sharding manager identity")
	data, err := c.bulkSyncByShardingManagerIdentity(ctxLogger, shardingManagerIdentity)
	if err != nil {
		return clusterConfigData, err
	}
	err = json.Unmarshal(data, &clusterConfigData)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to unmarshal cluster configuration")
		return clusterConfigData, err
	}
	return clusterConfigData, nil
}

func (r *registryClient) bulkSyncByShardingManagerIdentity(ctxLogger *logrus.Entry, shardingManagerIdentity string) ([]byte, error) {
	_, base, _, _ := runtime.Caller(0)
	absPath := filepath.Join(filepath.Dir(base), "/testdata/"+strings.TrimSpace(shardingManagerIdentity)+"-bulk.json")
	byteValue, err := os.ReadFile(absPath)
	if err != nil {
		ctxLogger.WithError(err).Error("failed perform bulk sync for cluster configuration from registry")
		return []byte(""), err
	}
	return byteValue, nil
}

func (c *registryClient) GetIdentitiesByCluster(ctx context.Context, clusterName string) (IdentityConfig, error) {
	var (
		identityConfig IdentityConfig
		tid            = uuid.NewString()
		ctxLogger      = log.WithFields(log.Fields{
			"clusterName": clusterName,
			"tid":         tid,
		})
	)
	ctxLogger.Infof("Get cluster configuration for provided sharding manager identity")
	data, err := c.getIdentitiesByCluster(ctxLogger, clusterName)
	if err != nil {
		return identityConfig, err
	}
	err = json.Unmarshal(data, &identityConfig)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to unmarshal cluster configuration")
		return identityConfig, err
	}
	return identityConfig, nil
}

func (r *registryClient) getIdentitiesByCluster(ctxLogger *logrus.Entry, clusterName string) ([]byte, error) {
	var (
		data = []byte("")
		err  error
	)
	_, base, _, _ := runtime.Caller(0)
	absPath := filepath.Join(filepath.Dir(base), "/testdata/"+clusterName+".json")
	data, err = os.ReadFile(absPath)
	if err != nil {
		ctxLogger.WithError(err).Error("failed to get cluster configuration from registry")
		return data, err
	}
	return data, nil
}
