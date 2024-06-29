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
	shardingHandler controller.ShardInteface
	mux             *http.ServeMux
	options         *options
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
	//initialize sharding manager with bootstrap configuration
	smConfig, err := manager.BootstrapConfiguration(ctx, params)
	if err != nil {
		log.Fatalf("failed to initialize sharding manager")
	}

	//initialize shard handler
	shardingHandler := controller.NewShardHandler(smConfig, params)
	err = shardingHandler.HandleLoadDistribution(ctx)
	if err != nil {
		log.Fatalf("error occurred while distributing load among operators: %v", err)
	}
	log.Println("sharding manager initialized")
	httpServer := &server{
		shardingHandler: shardingHandler,
		options:         createOptions(opts...),
		mux:             http.NewServeMux(),
	}
	httpServer.mux.HandleFunc(livenessPath, httpServer.livenessHandler)
	httpServer.mux.HandleFunc(readinessPath, httpServer.readinessHandler)
	return httpServer, nil
}

func (s *server) Listen(port string) error {
	return http.ListenAndServe(":"+port, s.mux)
}

func (s *server) livenessHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(200)
	_, err := responseWriter.Write([]byte(fmt.Sprintf("OK\n")))
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
	_, err := responseWriter.Write([]byte(fmt.Sprintf("OK\n")))
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
