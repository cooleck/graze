package swagger_ui

import (
	"net/http"

	"github.com/rakyll/statik/fs"
)

type Handler struct {
	fileServer http.Handler
}

func NewHandler() (*Handler, error) {
	statikFS, err := fs.NewWithNamespace(SwaggerUi)
	if err != nil {
		return nil, err
	}

	return &Handler{
		fileServer: http.FileServer(statikFS),
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fileServer.ServeHTTP(w, r)
}
