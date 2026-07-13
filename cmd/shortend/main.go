// Command shortend serves the shorten-service: creating a short code for a
// target URL over HTTP.
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/wmneco/wegweiser/internal/platform"
	"github.com/wmneco/wegweiser/internal/shorten"
)

func main() {
	cfg := platform.LoadConfig()
	logger := platform.NewLogger(cfg.LogLevel)

	core := shorten.NewCore()
	handler := shorten.NewHandler(core)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", platform.HealthHandler)
	mux.Handle("POST /api/shorten", handler)

	if err := platform.Run(context.Background(), logger, cfg, mux); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
