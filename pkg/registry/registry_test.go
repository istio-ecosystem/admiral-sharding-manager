package registry

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func getExpectedClusterConfiguration() ShardClusterConfig {
	cluster1 := ClusterConfig{
		Name:     "cluster1",
		Locality: "us-west-2",
		Metadata: clusterMetadata{},
	}
	cluster2 := ClusterConfig{
		Name:     "cluster2",
		Locality: "us-east-2",
		Metadata: clusterMetadata{},
	}

	var clusterConfigs = []ClusterConfig{}

	clusterConfigs = append(clusterConfigs, cluster1)
	clusterConfigs = append(clusterConfigs, cluster2)

	shardClusterConfig := ShardClusterConfig{
		Clusters:        clusterConfigs,
		LastUpdatedTime: "2024-06-20T18:25:43.511Z",
		ResourceVersion: "1.2.3",
	}

	return shardClusterConfig
}

func getExpectedBulkClusterConfiguration() ShardClusterConfig {
	cluster1 := ClusterConfig{
		Name:     "cluster1",
		Locality: "us-west-2",
		Metadata: clusterMetadata{},
		IdentityConfig: identityConfig{
			ClusterName: "cluster1",
			AssetList: []assetList{{
				Name:             "identity1",
				SourceAsset:      true,
				DestinationAsset: false,
			}},
		},
	}
	cluster2 := ClusterConfig{
		Name:     "cluster2",
		Locality: "us-east-2",
		Metadata: clusterMetadata{},
	}

	var clusterConfigs = []ClusterConfig{}

	clusterConfigs = append(clusterConfigs, cluster1)
	clusterConfigs = append(clusterConfigs, cluster2)

	shardClusterConfig := ShardClusterConfig{
		Clusters:        clusterConfigs,
		LastUpdatedTime: "2024-06-20T18:25:43.511Z",
		ResourceVersion: "1.2.3",
	}

	return shardClusterConfig
}

func getExpectedIdentityConfiguration() identityConfig {
	expectedIdentityConfig := identityConfig{
		ClusterName: "cluster1",
		AssetList: []assetList{{
			Name:             "identity1",
			SourceAsset:      true,
			DestinationAsset: false,
		}},
	}

	return expectedIdentityConfig
}

func TestParsingClusterConfig(t *testing.T) {
	expectedShardClusterConfig := getExpectedClusterConfiguration()
	testCases := []struct {
		name               string
		shardClusterConfig ShardClusterConfig
	}{
		{
			name: "Given a JSON cluster configuration, " +
				"When it is unmarshalled, " +
				"Then it should be read into the ClusterConfig struct",
			shardClusterConfig: expectedShardClusterConfig,
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			jsonResult, err := json.MarshalIndent(c.shardClusterConfig, "", "    ")
			if err != nil {
				t.Errorf("while marshaling ClusterConfig struct into JSON, got error: %s", err)
			}

			var unmarshalledClusterConfig ShardClusterConfig
			err = json.Unmarshal(jsonResult, &unmarshalledClusterConfig)
			if err != nil {
				t.Errorf("while unmarshaling JSON into ClusterConfig struct, got error: %s", err)
			}

			if !reflect.DeepEqual(unmarshalledClusterConfig, c.shardClusterConfig) {
				t.Errorf("unmarshalled ClusterConfig does not match with expected ClusterConfig, actual - %v, expected - %v", unmarshalledClusterConfig, expectedShardClusterConfig)
			}
		})
	}
}

func TestParsingIdentityConfig(t *testing.T) {
	expectedIdentityConfig := getExpectedIdentityConfiguration()
	testCases := []struct {
		name           string
		identityConfig identityConfig
	}{
		{
			name: "Given a JSON identity configuration, " +
				"When it is unmarshalled, " +
				"Then the config should be read into the identityConfig struct",
			identityConfig: expectedIdentityConfig,
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			jsonResult, err := json.MarshalIndent(c.identityConfig, "", "    ")
			if err != nil {
				t.Errorf("while marshaling identityConfig struct into JSON, got error: %s", err)
			}

			var unmarshalledClusterConfig identityConfig
			err = json.Unmarshal(jsonResult, &unmarshalledClusterConfig)
			if err != nil {
				t.Errorf("while unmarshaling JSON into identityConfig struct, got error: %s", err)
			}

			if !reflect.DeepEqual(unmarshalledClusterConfig, c.identityConfig) {
				t.Errorf("unmarshalled identityConfig does not match with expected identityConfig")
			}
		})
	}
}

func TestGetClustersByShardingManagerIdentity(t *testing.T) {
	expectedClusterConfig := getExpectedClusterConfiguration()
	registryClient := NewRegistryClient(WithEndpoint("endpoint"))
	testCases := []struct {
		name                  string
		expectedClusterConfig ShardClusterConfig
		expectedError         any
		smIdentity            string
		rc                    RegistryConfigInterface
	}{
		{
			name: "Given an sharding manager identity, " +
				"When GetClustersByShardingManagerIdentity is called, " +
				"Then actual config should match the expected cluster configuration",
			expectedClusterConfig: expectedClusterConfig,
			expectedError:         nil,
			smIdentity:            "test-shard-identity",
			rc:                    registryClient,
		},
		{
			name: "Given a non-existing sharding manager identity, " +
				"When GetClustersByShardingManagerIdentity is called, " +
				"Then there should be non nil error",
			expectedClusterConfig: expectedClusterConfig,
			expectedError:         errors.New(""),
			smIdentity:            "non-existing-shard-identity",
			rc:                    registryClient,
		},
		{
			name: "Given a existing sharding manager identity, " +
				"When GetClustersByShardingManagerIdentity is called and registry returns mis-configured json response, " +
				"Then there should be non nil error",
			expectedClusterConfig: expectedClusterConfig,
			expectedError:         errors.New(""),
			smIdentity:            "error-test-shard-identity",
			rc:                    registryClient,
		},
		{
			name: "Given a sharding manager identity and registry client is not initialized, " +
				"When GetClustersByShardingManagerIdentity is called, " +
				"Then there should be non nil error",
			expectedClusterConfig: expectedClusterConfig,
			expectedError:         errors.New(""),
			smIdentity:            "error-test-shard-identity",
			rc:                    NewRegistryClient(WithEndpoint("")),
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			actualClusterConfiguration, err := c.rc.GetClustersByShardingManagerIdentity(ctx, c.smIdentity)
			if err != nil && c.expectedError == nil {
				t.Errorf("error while getting cluster configuration with sharding manager identity, error: %v", err)
			}

			if err != nil && c.expectedError != nil && !errors.As(err, &c.expectedError) {
				t.Errorf("failed to get correct error: %v, instead got error: %v", c.expectedError, err)
			}

			if err == nil {
				if !cmp.Equal(actualClusterConfiguration, c.expectedClusterConfig) {
					t.Errorf("actual and expected ClusterConfig do not match. actual configuration : %v, expected configuration: %v", actualClusterConfiguration, c.expectedClusterConfig)
					t.Errorf(cmp.Diff(actualClusterConfiguration, c.expectedClusterConfig))
				}
			}
		})
	}
}

func TestBulkSyncByShardingManagerIdentity(t *testing.T) {
	expectedClusterConfig := getExpectedBulkClusterConfiguration()
	registryClient := NewRegistryClient(WithEndpoint("endpoint"))
	testCases := []struct {
		name                  string
		expectedClusterConfig ShardClusterConfig
		expectedError         any
		smIdentity            string
		rc                    RegistryConfigInterface
	}{
		{
			name: "Given an sharding manager identity, " +
				"When BulkSyncByShardingManagerIdentity is called, " +
				"Then actual config should match the expected cluster configuration",
			expectedClusterConfig: expectedClusterConfig,
			expectedError:         nil,
			smIdentity:            "test-shard-identity",
			rc:                    registryClient,
		},
		{
			name: "Given a non-existing sharding manager identity, " +
				"When BulkSyncByShardingManagerIdentity is called, " +
				"Then there should be non nil error",
			expectedClusterConfig: expectedClusterConfig,
			expectedError:         errors.New(""),
			smIdentity:            "non-existing-shard-identity",
			rc:                    registryClient,
		},
		{
			name: "Given a existing sharding manager identity, " +
				"When BulkSyncByShardingManagerIdentity is called and registry returns mis-configured json response, " +
				"Then there should be non nil error",
			expectedClusterConfig: expectedClusterConfig,
			expectedError:         errors.New(""),
			smIdentity:            "error-test-shard-identity",
			rc:                    registryClient,
		},
		{
			name: "Given a sharding manager identity and registry client is not initialized, " +
				"When BulkSyncByShardingManagerIdentity is called, " +
				"Then there should be non nil error",
			expectedClusterConfig: expectedClusterConfig,
			expectedError:         errors.New(""),
			smIdentity:            "error-test-shard-identity",
			rc:                    NewRegistryClient(WithEndpoint("")),
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			actualClusterConfiguration, err := c.rc.BulkSyncByShardingManagerIdentity(ctx, c.smIdentity)
			if err != nil && c.expectedError == nil {
				t.Errorf("error while getting cluster configuration with sharding manager identity, error: %v", err)
			}

			if err != nil && c.expectedError != nil && !errors.As(err, &c.expectedError) {
				t.Errorf("failed to get correct error: %v, instead got error: %v", c.expectedError, err)
			}

			if err == nil {
				if !cmp.Equal(actualClusterConfiguration, c.expectedClusterConfig) {
					t.Errorf("actual and expected ClusterConfig do not match. actual configuration : %v, expected configuration: %v", actualClusterConfiguration, c.expectedClusterConfig)
					t.Errorf(cmp.Diff(actualClusterConfiguration, c.expectedClusterConfig))
				}
			}
		})
	}
}

func TestGetIdentitiesByCluster(t *testing.T) {
	expectedIdentityConfig := getExpectedIdentityConfiguration()
	registryClient := NewRegistryClient(WithEndpoint("endpoint"))

	testCases := []struct {
		name                   string
		expectedIdentityConfig identityConfig
		expectedError          any
		clusterName            string
		rc                     RegistryConfigInterface
	}{
		{
			name: "Given a cluster name, " +
				"When GetIdentitiesByCluster is called, " +
				"Then actual config should match the expected identity configuration",
			expectedIdentityConfig: expectedIdentityConfig,
			expectedError:          nil,
			clusterName:            "test-cluster-identity",
			rc:                     registryClient,
		},
		{
			name: "Given a non-existing cluster name, " +
				"When GetIdentitiesByCluster is called, " +
				"Then there should be non nil error",
			expectedIdentityConfig: expectedIdentityConfig,
			expectedError:          errors.New(""),
			clusterName:            "non-existing-cluster",
			rc:                     registryClient,
		},
		{
			name: "Given a non-existing cluster name, " +
				"When GetIdentitiesByCluster is called and registry returns mis-configured identity json response, " +
				"Then there should be non nil error",
			expectedIdentityConfig: expectedIdentityConfig,
			expectedError:          errors.New(""),
			clusterName:            "error-test-cluster-identity",
			rc:                     registryClient,
		},
		{
			name: "Given a non-existing cluster name, " +
				"When GetIdentitiesByCluster is called and registry client is not initialized, " +
				"Then there should be non nil error",
			expectedIdentityConfig: expectedIdentityConfig,
			expectedError:          errors.New(""),
			clusterName:            "error-test-cluster-identity",
			rc:                     NewRegistryClient(WithEndpoint("")),
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			actualClusterConfiguration, err := c.rc.GetIdentitiesByCluster(ctx, c.clusterName)
			if err != nil && c.expectedError == nil {
				t.Errorf("error while getting cluster configuration with sharding manager identity, error: %v", err)
			}

			if err != nil && c.expectedError != nil && !errors.As(err, &c.expectedError) {
				t.Errorf("failed to get correct error: %v, instead got error: %v", c.expectedError, err)
			}

			if err == nil {
				if !cmp.Equal(actualClusterConfiguration, c.expectedIdentityConfig) {
					t.Errorf("actual and expected ClusterConfig do not match. actual configuration : %v, expected configuration: %v", actualClusterConfiguration, c.expectedIdentityConfig)
					t.Errorf(cmp.Diff(actualClusterConfiguration, c.expectedIdentityConfig))
				}
			}
		})
	}
}
