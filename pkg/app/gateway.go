package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func (a *App) initGateway(cfg *Config) {
	ctx := context.Background()

	mux := runtime.NewServeMux()

	for _, service := range cfg.Services {
		if err := service.RegisterGateway(ctx, mux); err != nil {
			logger.Fatal("failed to register gateway", zap.Error(err))
		}
	}

	toolsOrigin := fmt.Sprintf("http://localhost:%d", cfg.ToolsPort)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			toolsOrigin,
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"ResponseType",
			"api_key",
			"Authorization",
		},
	})

	handler := c.Handler(mux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.GatewayPort),
		Handler: handler,
	}

	a.gatewayServer = httpServer
}

func (a *App) runGateway() error {
	if err := a.gatewayServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
