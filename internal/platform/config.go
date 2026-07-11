// Package platform provides shared HTTP service runtime mechanics — server
// lifecycle, middleware, health checks, error responses, logging, and
// configuration — used by every Wegweiser service binary. It contains no
// URL-shortener domain knowledge.
package platform

import (
	"os"
	"time"
)

// Environment variable names read by LoadConfig.
const (
	EnvAddr            = "WEGWEISER_ADDR"
	EnvLogLevel        = "WEGWEISER_LOG_LEVEL"
	EnvShutdownTimeout = "WEGWEISER_SHUTDOWN_TIMEOUT"
)

// Default configuration values used when the corresponding environment
// variable is unset or invalid.
const (
	DefaultAddr            = ":8080"
	DefaultLogLevel        = "info"
	DefaultShutdownTimeout = 5 * time.Second
)

// Config holds runtime configuration shared by every service binary.
type Config struct {
	// Addr is the address the HTTP server listens on, e.g. ":8080".
	Addr string
	// LogLevel is the minimum slog level to emit: debug, info, warn, or error.
	LogLevel string
	// ShutdownTimeout bounds how long graceful shutdown waits for in-flight
	// requests to finish before the process exits.
	ShutdownTimeout time.Duration
}

// LoadConfig reads Config from the environment, falling back to documented
// defaults for any value that is unset or fails to parse.
func LoadConfig() Config {
	return Config{
		Addr:            stringEnv(EnvAddr, DefaultAddr),
		LogLevel:        stringEnv(EnvLogLevel, DefaultLogLevel),
		ShutdownTimeout: durationEnv(EnvShutdownTimeout, DefaultShutdownTimeout),
	}
}

func stringEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func durationEnv(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
