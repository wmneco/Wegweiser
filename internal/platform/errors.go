package platform

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the consistent JSON shape used for every error response
// across services.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteError writes status and a JSON body of shape ErrorResponse to w. It is
// the single helper every handler and middleware uses to report errors, so
// all error responses share the same shape.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Code: code, Message: message})
}
