// Command redirectd serves the redirect-service: resolving a short code to
// its target URL over HTTP.
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/wmneco/wegweiser/internal/platform"
	"github.com/wmneco/wegweiser/internal/redirect"
)

func main() {
	cfg := platform.LoadConfig()
	logger := platform.NewLogger(cfg.LogLevel)

	core := redirect.NewCore()
	handler := redirect.NewHandler(core)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", platform.HealthHandler)
	mux.Handle("GET /{code}", handler)

	if err := platform.Run(context.Background(), logger, cfg, mux); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
