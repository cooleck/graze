package swagger_json

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/spec"
	"net"
	"net/http"

	"github.com/peterbourgon/mergemap"
)

type Handler struct {
	swaggerDescMerged []byte
	gatewayPort       int
}

func NewHandler(gatewayPort int, swaggerDescriptions [][]byte) (*Handler, error) {
	m1 := map[string]interface{}{}
	m2 := map[string]interface{}{}
	for _, swaggerDescription := range swaggerDescriptions {
		err := json.Unmarshal(swaggerDescription, &m1)
		if err != nil {
			return nil, err
		}
		m2 = mergemap.Merge(m2, m1)
	}

	swaggerDescMerged, err := json.Marshal(m2)
	if err != nil {
		return nil, err
	}

	return &Handler{
		swaggerDescMerged: swaggerDescMerged,
		gatewayPort:       gatewayPort,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var sw spec.Swagger
	if err = json.Unmarshal(h.swaggerDescMerged, &sw); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sw.Host = fmt.Sprintf("%s:%d", host, h.gatewayPort)
	b, err := json.Marshal(sw)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b)
}
