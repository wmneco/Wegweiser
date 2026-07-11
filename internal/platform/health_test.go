package platform_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wmneco/wegweiser/internal/platform"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	platform.HealthHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
