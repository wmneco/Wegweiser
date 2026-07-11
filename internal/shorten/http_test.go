package shorten_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wmneco/wegweiser/internal/platform"
	"github.com/wmneco/wegweiser/internal/shorten"
)

func TestHandler_PlaceholderReturns501(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("POST /api/shorten", shorten.NewHandler(shorten.NewCore()))

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotImplemented, rec.Code)

	var body platform.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.NotEmpty(t, body.Code)
	assert.NotEmpty(t, body.Message)
}
