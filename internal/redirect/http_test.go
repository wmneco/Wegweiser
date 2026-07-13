package redirect_test

import (
	"context"
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

type successResolver struct{}

func (successResolver) Resolve(context.Context, string) (string, error) {
	return "https://example.com", nil
}

func TestHandler_PlaceholderReturns501_EvenWhenResolverSucceeds(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("GET /{code}", redirect.NewHandler(successResolver{}))

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotImplemented, rec.Code)
}
