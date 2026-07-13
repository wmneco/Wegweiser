package platform

import "net/http"

// HealthHandler responds to liveness checks with 200 and a small JSON body,
// without depending on any downstream store or domain logic.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
