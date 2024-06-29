/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/controller"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/manager"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/monitoring"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

// discoveryCmd represents the discovery command
var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
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

		//initialize sharding manager with bootstrap configuration
		smConfig, err := manager.InitializeShardingManager(ctx, &smParams)
		if err != nil {
			log.Fatalf("failed to initialize sharding manager")
		}

		//initialize shard handler
		shardingHandler := controller.NewShardHandler(smConfig, &smParams)
		err = shardingHandler.HandleLoadDistribution(ctx)
		if err != nil {
			log.Fatalf("error occurred while distributing load among operators: %v", err)
		}
	},
}

func init() {
	discoveryCmd.PersistentFlags().StringVar(&smParams.KubeconfigPath, "kube_config", "", "Use a Kubernetes configuration file instead of in-cluster configuration")
	//defines the identity of sharding manager instance - logical name for group of resources handled by an instance of sharding manager. This is used to initialize configuration from registry and as "admiral.io/shardingMangerIdentity" label value on shard crd
	discoveryCmd.Flags().StringVar(&smParams.ShardingManagerIdentity, "shard-identity", "devx", "Identity of the sharding manager instance, used to get configuration from registry and used as value for label \"admiral.io/shardingMangerIdentity\" on shard crd ")
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
	newServer, err := server.NewServer()
	if err != nil {
		log.Fatalf("failed to instantiate a server: %v", err)
	}
	err = newServer.Listen(model.PortNumber)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Printf("started server on port %s", model.PortNumber)
}

// intialize monitoring
func initializeMonitoring() {
	err := monitoring.InitializeMonitoring()
	if err != nil {
		log.Fatalf("failed to initialize monitoring: %v", err)
	}
}

func startConfigDiscovery() {
}
