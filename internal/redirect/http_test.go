package redirect_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wmneco/wegweiser/internal/platform"
	"github.com/wmneco/wegweiser/internal/redirect"
)

func TestHandler_PlaceholderReturns501(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("GET /{code}", redirect.NewHandler(redirect.NewCore()))

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotImplemented, rec.Code)

	var body platform.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.NotEmpty(t, body.Code)
	assert.NotEmpty(t, body.Message)
}
