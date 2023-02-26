package requestdump

import (
	"context"
	"net/http"
	"net/http/httputil"

	"github.com/Bhargav-InfraCloud/rate-limit-server/pkg/logs"
	"github.com/Bhargav-InfraCloud/rate-limit-server/pkg/middlewares"
)

const MiddlewareName middlewares.MiddlewareName = `request-dump`

type Handler struct {
	logger logs.Logger
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{
		logger: logs.FromContext(ctx),
	}
}

func (h *Handler) RequestDump(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := h.logger.With().Stringer("middleware", MiddlewareName).
			Str("path", r.URL.Path).Logger()

		logger.Info().Msg("New request through middleware")

		defer next.ServeHTTP(w, r)

		requestBytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to dump request")

			return
		}

		logger.Debug().Bytes("dump", requestBytes).Msg("Log request dump")
	})
}
