package redirect

import (
	"net/http"

	"github.com/wmneco/wegweiser/internal/platform"
)

// Handler is the HTTP driving adapter for GET /{code}. It invokes the
// inbound Resolver port rather than embedding behaviour directly.
type Handler struct {
	resolver Resolver
}

// NewHandler builds a Handler backed by resolver.
func NewHandler(resolver Resolver) *Handler {
	return &Handler{resolver: resolver}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	_, _ = h.resolver.Resolve(r.Context(), code)
	platform.WriteError(w, http.StatusNotImplemented, "not_implemented", "redirect is not yet implemented")
}
