package app

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"

	swagger_json "github.com/cooleck/graze/internal/pkg/swagger/swagger-json"
	swagger_ui "github.com/cooleck/graze/internal/pkg/swagger/swagger-ui"
)

func (a *App) initTools(cfg *Config) {
	mux := http.NewServeMux()

	swaggerDescriptions := make([][]byte, 0, len(cfg.Services))
	for _, service := range cfg.Services {
		swaggerDescriptions = append(swaggerDescriptions, service.GetSwaggerDesc())
	}

	swaggerJSONHandler, err := swagger_json.NewHandler(a.gatewayPort, swaggerDescriptions)
	if err != nil {
		logger.Fatal("failed to create Swagger JSON handler", zap.Error(err))
	}
	mux.Handle("/docs/swagger.json", swaggerJSONHandler)

	swaggerUIHandler, err := swagger_ui.NewHandler()
	if err != nil {
		logger.Fatal("failed to create Swagger UI handler", zap.Error(err))
	}
	mux.Handle("/docs/", http.StripPrefix("/docs/", swaggerUIHandler))

	a.toolsServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.toolsPort),
		Handler: mux,
	}
}

func (a *App) runTools() error {
	return a.toolsServer.ListenAndServe()
}
