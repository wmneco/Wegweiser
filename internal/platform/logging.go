package platform

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger builds a structured (JSON) slog.Logger honoring the given level
// name (debug, info, warn, error — case-insensitive). Unrecognized levels
// fall back to info.
func NewLogger(level string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseLevel(level)})
	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
