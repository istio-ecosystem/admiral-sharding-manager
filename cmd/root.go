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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	smParams := model.ShardingManagerParams{}

	rootCmd.PersistentFlags().StringVar(&smParams.KubeconfigPath, "kube_config", "", "Use a Kubernetes configuration file instead of in-cluster configuration")
	rootCmd.Flags().StringVar(&smParams.ShardingManagerIdentity, "sharding-manager-identity", "devx", "Identity of the sharding manager instance")
	rootCmd.Flags().StringVar(&smParams.OperatorIdentityLabel, "shard-workload-identity-label", "sharding.manager.io/shard-workload-identity", "label used to specify identity of workload for whome shard profile is defined")
	rootCmd.Flags().StringVar(&smParams.ShardIdentityLabel, "shard-instance-identity-label", "sharding.manager.io/shard-instance-identity", "label used to specify identity of sharding manager instance")
	rootCmd.Flags().StringVar(&smParams.ShardNamespace, "shard-namespace", "shard-namespace", "Namespace used to create sharding resources")

	//fetch bootstrap configuration from registry
	ctx := context.Background()

	smConfig, err := manager.InitializeShardingManager(ctx, &smParams)
	if err != nil {
		logrus.Error("failed to initialize sharding manager")
	}

	controller.NewShardHandler(smConfig, smParams.ShardNamespace)
}
