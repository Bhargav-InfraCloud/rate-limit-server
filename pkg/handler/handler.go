package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	pkgDB "github.com/Bhargav-InfraCloud/rate-limit-server/pkg/database"
	pkgLogs "github.com/Bhargav-InfraCloud/rate-limit-server/pkg/logs"
	pkgSvc "github.com/Bhargav-InfraCloud/rate-limit-server/pkg/service"
	pkgSvcErr "github.com/Bhargav-InfraCloud/rate-limit-server/pkg/service/errors"
	mux "github.com/vaguecoder/gorilla-mux"
)

const (
	idKey    key = `id`
	countKey key = `count`
	resetKey key = `reset`

	statusOK = `ok`
	trueStr  = `true`
)

type RateLimitHandler struct {
	db pkgDB.RateLimiterDB
}

func NewRateLimitHandler(db pkgDB.RateLimiterDB) *RateLimitHandler {
	return &RateLimitHandler{
		db: db,
	}
}

func (h *RateLimitHandler) HandlerFunc(ctx context.Context) http.HandlerFunc {
	var (
		err          error
		servErr      *pkgSvcErr.Error
		count        int
		id, countStr string
		respBody     response
		resetCount   bool
		logs         []string
		muxVars      map[string]string

		logger = pkgLogs.FromContext(ctx)
	)

	return func(w http.ResponseWriter, r *http.Request) {
		logs = []string{}

		muxVars = mux.Vars(r)
		id = muxVars[idKey.String()]

		countStr = r.Header.Get(countKey.String())
		if countStr != "" {
			count, err = strconv.Atoi(countStr)
			if err != nil {
				logs = []string{
					"unparsable count value in header",
					err.Error(),
				}
				h.writeErr(ctx, w, pkgSvcErr.InvalidInputsError, logs)

				return
			} else if count < 0 {
				logs = []string{
					fmt.Sprintf("count %q is not an unsigned integer", count),
				}
				h.writeErr(ctx, w, pkgSvcErr.InvalidInputsError, logs)

				return
			}
		}

		resetCount = r.Header.Get(resetKey.String()) == trueStr

		logger = pkgLogs.FromRawLogger(logger.With().Str("id", id).
			Int("count", count).Bool("reset", resetCount).Logger())

		servErr = h.db.Add(id, uint(count), resetCount)
		if servErr != nil {
			logs = append(logs, fmt.Sprintf("failed to add id %q", id))
			h.writeErr(ctx, w, servErr, logs)

			return
		}

		respBody = response{
			Status: http.StatusText(http.StatusOK),
		}

		h.writeResponse(ctx, w, respBody)
	}
}

func (h *RateLimitHandler) writeErr(ctx context.Context, w http.ResponseWriter, servErr *pkgSvcErr.Error, logs []string) {
	var (
		err    error
		logger = pkgLogs.FromContext(ctx)
	)

	if servErr != nil {
		logger = pkgLogs.FromRawLogger(logger.With().Str("service-code", servErr.Code).
			Int("status-code", servErr.StatusCode()).
			Strs("related-logs", servErr.RelatedLogs).Logger())
	}

	logger.Info().Err(servErr).Strs("logs", logs).Msg("Writing error to response writer")

	if servErr == nil {
		logger.Info().Msg("Service error is nil")

		return
	}

	servErr = pkgSvcErr.CopyOfServiceError(servErr)

	if logs != nil {
		servErr.RelatedLogs = append(servErr.RelatedLogs, logs...)
	}

	err = pkgSvc.WriteErrorResponse(w, *servErr)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write error response")
	}

	logger.Info().Msg("Successfully written error to response")
}

func (h *RateLimitHandler) writeResponse(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	var (
		err error

		logger = pkgLogs.FromContext(ctx).With().Interface("payload", payload).Logger()
	)

	err = pkgSvc.WriteResponse(w, payload)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write response")

		return
	}

	logger.Info().Msg("Successfully written response")
}
