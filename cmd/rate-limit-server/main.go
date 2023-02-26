package main

import (
	"context"
	"net/http"
	"os"

	"github.com/Bhargav-InfraCloud/rate-limit-server/pkg/database"
	"github.com/Bhargav-InfraCloud/rate-limit-server/pkg/handler"
	"github.com/Bhargav-InfraCloud/rate-limit-server/pkg/logs"
	"github.com/Bhargav-InfraCloud/rate-limit-server/pkg/middlewares/requestdump"
	mux "github.com/vaguecoder/gorilla-mux"
)

const (
	rateLimitHandlerPath = `/id/{id}`
	addr                 = `:8080`
)

func main() {
	var (
		err          error
		logger       logs.Logger
		mws          []mux.MiddlewareFunc
		handlerPaths []string

		ctx    = context.Background()
		dbOps  = database.NewDB()
		router = mux.NewRouter()
	)

	ctx, logger = logs.NewLogger(ctx, os.Stdout, logs.LevelDebug)

	mws = []mux.MiddlewareFunc{
		requestdump.NewHandler(ctx).RequestDump,
	}

	router.Use(mws...)

	rateLimitHandler := handler.NewRateLimitHandler(dbOps).HandlerFunc(ctx)

	router.Handle(rateLimitHandlerPath, rateLimitHandler).Methods(http.MethodGet)
	handlerPaths = append(handlerPaths, rateLimitHandlerPath)

	logger.Info().Strs("paths", handlerPaths).Str("address", addr).Msg("Serving on address")

	if err = http.ListenAndServe(addr, router); err != nil {
		logger.Error().Err(err).Msg("Failure at serving")
	}
}
