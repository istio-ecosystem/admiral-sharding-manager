package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/controller"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/manager"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/monitoring"
	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/registry"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

const (
	livenessPath  = "/liveness"
	readinessPath = "/readiness"
)

var (
	shardingManagerServerMeter   = monitoring.NewMeter("admiral_sharding_manager_server")
	shardingManagerRequestsTotal = monitoring.NewCounter(
		"requests_total",
		"total number of requests handled by sharding manager server",
		monitoring.WithMeter(shardingManagerServerMeter))
)

type server struct {
	mux     *http.ServeMux
	options *options
}

type options struct {
}

func createOptions(opts ...Options) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Options accepts a pointer to options. It is used
// to update the options by calling an array of functions
type Options func(*options)

func NewServer(ctx context.Context, params *model.ShardingManagerParams, opts ...Options) (*server, error) {
	client, err := initClients(params)
	if err != nil {
		return nil, fmt.Errorf("failed setting up clients: %v", err)
	}
	shardHandler := controller.NewShardHandler(client, params)
	shardingManager, err := manager.NewShardingManager(ctx, shardHandler, client, params.ShardingManagerIdentity)
	if err != nil {
		return nil, fmt.Errorf("error initializing sharding manager: %v", err)
	}
	err = shardingManager.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start sharding manager")
	}

	httpServer := &server{
		options: createOptions(opts...),
		mux:     http.NewServeMux(),
	}
	httpServer.mux.HandleFunc(livenessPath, httpServer.livenessHandler)
	httpServer.mux.HandleFunc(readinessPath, httpServer.readinessHandler)
	return httpServer, nil
}

func initClients(params *model.ShardingManagerParams) (model.Clients, error) {
	var client model.Clients
	var kubeClient manager.LoadKubeClient = &manager.KubeClient{}
	admiralAPIClient, err := kubeClient.LoadAdmiralApiClientFromPath(params.KubeconfigPath)
	if err != nil {
		return client, fmt.Errorf("failed to initialize admiral api client")
	}
	client.AdmiralClient = admiralAPIClient
	client.RegistryClient = registry.NewRegistryClient(registry.WithEndpoint(params.RegistryEndpoint))
	return client, nil
}

func (s *server) Listen(port string) error {
	return http.ListenAndServe(":"+port, s.mux)
}

func (s *server) livenessHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(200)
	_, err := responseWriter.Write([]byte(fmt.Sprintln("OK")))
	if err != nil {
		shardingManagerRequestsTotal.Increment(api.WithAttributes(
			attribute.Key("path").String(livenessPath),
			attribute.Key("code").String("503"),
		))
		log.Fatalf("failed to write response")
	}
	shardingManagerRequestsTotal.Increment(api.WithAttributes(
		attribute.Key("path").String(livenessPath),
		attribute.Key("code").String("200"),
	))
}

func (s *server) readinessHandler(responseWriter http.ResponseWriter, request *http.Request) {
	shardingManagerRequestsTotal.Increment(api.WithAttributes(
		attribute.Key("path").String(readinessPath),
	))
	responseWriter.WriteHeader(200)
	_, err := responseWriter.Write([]byte(fmt.Sprintln("OK")))
	if err != nil {
		shardingManagerRequestsTotal.Increment(api.WithAttributes(
			attribute.Key("path").String(readinessPath),
			attribute.Key("code").String("503"),
		))
		log.Fatalf("failed to write response")
	}
	shardingManagerRequestsTotal.Increment(api.WithAttributes(
		attribute.Key("path").String(readinessPath),
		attribute.Key("code").String("200"),
	))
}
