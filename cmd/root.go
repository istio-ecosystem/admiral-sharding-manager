/*
Copyright © 2024 Intuit Inc.
*/
package cmd

import (
	"context"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/controller"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/manager"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "admiral-sharding-manager",
	Short: "Sharding manager distributes load among admiral operators",
	Long:  "Sharding manager distributes load among admiral operators",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	smParams := model.ShardingManagerParams{}

	rootCmd.PersistentFlags().StringVar(&smParams.KubeconfigPath, "kube_config", "", "Use a Kubernetes configuration file instead of in-cluster configuration")
	//defines the identity of sharding manager instance - logical name for group of resources handled by an instance of sharding manager. This is used to initialize configuration for resources to be handled
	rootCmd.Flags().StringVar(&smParams.ShardingManagerIdentity, "shard-identity", "devx", "Identity of the sharding manager instance")
	//operator identity label which will be set on the shard crd. Using this label value operator will filter the shard it needs to monitor
	rootCmd.Flags().StringVar(&smParams.OperatorIdentityLabel, "operator-identity-label", "sharding.manager.io/operator-identity", "label used to specify identity of operator for which shard profile is defined")
	//shard identity label, value of which will be sharding manager instance identity. This is added on the shard crd to identity which sharding manager instance managed this shard crd.
	rootCmd.Flags().StringVar(&smParams.ShardIdentityLabel, "shard-identity-label", "sharding.manager.io/shard-identity", "label used to specify identity of sharding manager instance")
	//shard namespace defines the namspace in which sharding manager should drop in shard crds
	rootCmd.Flags().StringVar(&smParams.ShardNamespace, "shard-namespace", "shard-namespace", "Namespace used to create sharding resources")

	//fetch bootstrap configuration from registry
	ctx := context.Background()

	smConfig, err := manager.InitializeShardingManager(ctx, &smParams)
	if err != nil {
		logrus.Error("failed to initialize sharding manager")
	}

	controller.NewShardHandler(smConfig, smParams.ShardNamespace)
}
