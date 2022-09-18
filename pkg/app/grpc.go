package app

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

const (
	tcpNetwork = "tcp"
)

func (a *App) initGRPC(cfg *Config) {
	a.grpcServer = grpc.NewServer(cfg.GRPCServerOpts...)
	for _, service := range cfg.Services {
		service.RegisterGRPC(a.grpcServer)
	}
}

func (a *App) runGRPC() error {
	grpcListener, err := net.Listen(tcpNetwork, fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(grpcListener)
}
