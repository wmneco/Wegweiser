package platform

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Run starts an HTTP server on cfg.Addr serving handler wrapped with the
// baseline Middleware, and blocks until ctx is done, the process receives
// SIGINT/SIGTERM, or the server fails. On signal, it stops accepting new
// connections and waits up to cfg.ShutdownTimeout for in-flight requests to
// finish before returning.
func Run(ctx context.Context, logger *slog.Logger, cfg Config, handler http.Handler) error {
	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: Middleware(logger, handler),
	}

	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-sigCtx.Done():
		logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
