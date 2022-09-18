package app

import (
	"context"
	"go.uber.org/zap"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const (
	grpcDefaultPort    = 8002
	gatewayDefaultPort = 8000
	toolsDefaultPort   = 8004
)

var (
	logger, _ = zap.NewProduction()
)

type Service interface {
	RegisterGRPC(s grpc.ServiceRegistrar)
	RegisterGateway(ctx context.Context, mux *runtime.ServeMux) error
	GetSwaggerDesc() []byte
}

type Config struct {
	GRPCPort       int
	GRPCServerOpts []grpc.ServerOption

	GatewayPort int

	ToolsPort int

	Services []Service
}

type App struct {
	grpcPort   int
	grpcServer *grpc.Server

	gatewayPort   int
	gatewayServer *http.Server

	toolsPort   int
	toolsServer *http.Server
}

func NewDefaultPorts(services []Service, grpcServerOpts ...grpc.ServerOption) *App {
	return New(&Config{
		GRPCPort:       grpcDefaultPort,
		GRPCServerOpts: grpcServerOpts,
		GatewayPort:    gatewayDefaultPort,
		ToolsPort:      toolsDefaultPort,
		Services:       services,
	})
}

func New(cfg *Config) *App {
	a := &App{
		grpcPort:    cfg.GRPCPort,
		gatewayPort: cfg.GatewayPort,
		toolsPort:   cfg.ToolsPort,
	}

	a.initGRPC(cfg)
	a.initGateway(cfg)
	a.initTools(cfg)

	return a
}

func (a *App) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return a.runGRPC()
	})

	group.Go(func() error {
		return a.runGateway()
	})

	group.Go(func() error {
		return a.runTools()
	})

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}
