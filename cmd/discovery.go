/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/monitoring"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	ctx = context.Background()
)

// discoveryCmd represents the discovery command
var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Discover configuration to distribute amongst admiral operators",
	Long:  `Discover configuration to distribute amongst admiral operators.`,
	Run: func(cmd *cobra.Command, args []string) {
		//initialize monitoring and start servers
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			Initialize(
				initializeMonitoring,
				startMetricsServer,
				startNewServer,
				startConfigDiscovery,
			)
			wg.Done()
		}()

		wg.Wait()
	},
}

func init() {
	discoveryCmd.PersistentFlags().StringVar(&smParams.KubeconfigPath, "kube_config", "", "Use a Kubernetes configuration file instead of in-cluster configuration")
	//defines the identity of sharding manager instance - logical name for group of resources handled by an instance of sharding manager. This is used to initialize configuration from registry and as "admiral.io/shardingMangerIdentity" label value on shard crd
	discoveryCmd.Flags().StringVar(&smParams.ShardingManagerIdentity, "shard-identity", "dev", "Identity of the sharding manager instance, used to get configuration from registry and used as value for label \"admiral.io/shardingMangerIdentity\" on shard crd ")
	//operator identity label which will be set on the shard crd. Using this label value operator will filter the shard it needs to monitor
	discoveryCmd.Flags().StringVar(&smParams.OperatorIdentityLabel, "operator-identity-label", "admiral.io/operatorIdentity", "label used to specify identity of operator for which shard profile is defined")
	//shard namespace defines the namspace in which sharding manager should drop in shard crds
	discoveryCmd.Flags().StringVar(&smParams.ShardNamespace, "shard-namespace", "shard-namespace", "Namespace used to create sharding resources")
	//registry endpoint
	discoveryCmd.Flags().StringVar(&smParams.RegistryEndpoint, "registry-endpoint", "", "Registry Service endpoint to get configuration for sharding manager")

	rootCmd.AddCommand(discoveryCmd)
}

func Initialize(funcs ...func()) {
	wg := new(sync.WaitGroup)
	wg.Add(len(funcs))
	for _, fn := range funcs {
		go func() {
			fn()
			wg.Done()
		}()
	}
	wg.Wait()
}

// initialize metrics service
func startMetricsServer() {
	http.Handle(model.MetricsPath, promhttp.Handler())
	err := http.ListenAndServe(":"+model.MetricsPort, nil)
	if err != nil {
		log.Fatalf("error serving http: %v", err)
	}
}

// initialize sharding manager server
func startNewServer() {

	newServer, err := server.NewServer(ctx, &smParams)
	if err != nil {
		log.Fatalf("failed to instantiate a server: %v", err)
	}
	err = newServer.Listen(model.PortNumber)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Printf("started server on port %s", model.PortNumber)
}

// initialize monitoring
func initializeMonitoring() {
	err := monitoring.InitializeMonitoring()
	if err != nil {
		log.Fatalf("failed to initialize monitoring: %v", err)
	}
}

func startConfigDiscovery() {
}
