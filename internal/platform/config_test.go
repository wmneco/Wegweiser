package platform_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wmneco/wegweiser/internal/platform"
)

func TestLoadConfig_DefaultsWhenUnset(t *testing.T) {
	t.Setenv(platform.EnvAddr, "")
	t.Setenv(platform.EnvLogLevel, "")
	t.Setenv(platform.EnvShutdownTimeout, "")

	cfg := platform.LoadConfig()

	assert.Equal(t, platform.DefaultAddr, cfg.Addr)
	assert.Equal(t, platform.DefaultLogLevel, cfg.LogLevel)
	assert.Equal(t, platform.DefaultShutdownTimeout, cfg.ShutdownTimeout)
}

func TestLoadConfig_EnvOverridesDefaults(t *testing.T) {
	t.Setenv(platform.EnvAddr, ":9090")
	t.Setenv(platform.EnvLogLevel, "debug")
	t.Setenv(platform.EnvShutdownTimeout, "10s")

	cfg := platform.LoadConfig()

	assert.Equal(t, ":9090", cfg.Addr)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
}

func TestLoadConfig_InvalidDurationFallsBackToDefault(t *testing.T) {
	t.Setenv(platform.EnvShutdownTimeout, "not-a-duration")

	cfg := platform.LoadConfig()

	assert.Equal(t, platform.DefaultShutdownTimeout, cfg.ShutdownTimeout)
}
