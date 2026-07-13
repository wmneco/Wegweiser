package shorten

import (
	"net/http"

	"github.com/wmneco/wegweiser/internal/platform"
)

// Handler is the HTTP driving adapter for POST /api/shorten. It invokes the
// inbound Shortener port rather than embedding behaviour directly.
type Handler struct {
	shortener Shortener
}

// NewHandler builds a Handler backed by shortener.
func NewHandler(shortener Shortener) *Handler {
	return &Handler{shortener: shortener}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = h.shortener.Shorten(r.Context(), "")
	platform.WriteError(w, http.StatusNotImplemented, "not_implemented", "shorten is not yet implemented")
}
