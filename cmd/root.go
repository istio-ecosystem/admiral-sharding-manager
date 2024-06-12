/*
Copyright © 2024 Intuit Inc.
*/
package cmd

import (
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/controller"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/manager"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sync"

	"context"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/server"
	"github.com/spf13/cobra"
	"os"
)

const (
	portNumber  = "8080"
	metricsPort = "9090"
	metricsPath = "/metrics"
)

// maintains sharding manger program arguments
var smParams = model.ShardingManagerParams{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "admiral-sharding-manager",
	Short: "Sharding manager distributes load among admiral operators",
	Long:  "Sharding manager distributes load among admiral operators",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		//initialize sharding manager with bootstrap configuration
		smConfig, err := manager.InitializeShardingManager(ctx, &smParams)
		if err != nil {
			log.Fatalf("failed to initialize sharding manager")
		}

		//initialize shard handler
		controller.NewShardHandler(smConfig, smParams.ShardNamespace)

		//initialize monitoring and start servers
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			Initialize(
				initializeMonitoring,
				startMetricsServer,
				startNewServer,
			)
			wg.Done()
		}()
		wg.Wait()
		log.Println("admiral sharding manager has been initialized")
	},
}

// initialize metrics service
func startMetricsServer() {
	http.Handle(metricsPath, promhttp.Handler())
	err := http.ListenAndServe(":"+metricsPort, nil)
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
	err = newServer.Listen(portNumber)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Printf("started server on port %s", portNumber)
}

// intialize monitoring
func initializeMonitoring() {
	err := monitoring.InitializeMonitoring()
	if err != nil {
		log.Fatalf("failed to initialize monitoring: %v", err)
	}
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// manage root command flags
func init() {
	rootCmd.PersistentFlags().StringVar(&smParams.KubeconfigPath, "kube_config", "", "Use a Kubernetes configuration file instead of in-cluster configuration")
	//defines the identity of sharding manager instance - logical name for group of resources handled by an instance of sharding manager. This is used to initialize configuration from registry and as "admiral.io/shardingMangerIdentity" label value on shard crd
	rootCmd.Flags().StringVar(&smParams.ShardingManagerIdentity, "shard-identity", "devx", "Identity of the sharding manager instance, used to get configuration from registry and used as value for label \"admiral.io/shardingMangerIdentity\" on shard crd ")
	//operator identity label which will be set on the shard crd. Using this label value operator will filter the shard it needs to monitor
	rootCmd.Flags().StringVar(&smParams.OperatorIdentityLabel, "operator-identity-label", "admiral.io/operatorIdentity", "label used to specify identity of operator for which shard profile is defined")
	//shard namespace defines the namspace in which sharding manager should drop in shard crds
	rootCmd.Flags().StringVar(&smParams.ShardNamespace, "shard-namespace", "shard-namespace", "Namespace used to create sharding resources")
	//registry endpoint
	rootCmd.Flags().StringVar(&smParams.RegistryEndpoint, "registry-endpoint", "", "Registry Service endpoint to get configuration for sharding manager")

}
