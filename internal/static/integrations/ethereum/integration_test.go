package ethereum

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Disabled(t *testing.T) {
	cfg := &Config{
		Enabled: false,
	}
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()

	integration := New(logger, cfg, reg)
	err := integration.Run(context.Background())
	require.NoError(t, err)
}

func TestIntegration_ExecutionClient(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Execution: ExecutionConfig{
			Enabled: true,
			URL:     "http://localhost:8545", // This should be a mock server in real tests
			Modules: []string{"eth", "net"},
		},
	}
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()

	integration := New(logger, cfg, reg)

	// Start the integration
	err := integration.Run(context.Background())
	require.Error(t, err) // Should error because we can't connect to the client
	assert.Contains(t, err.Error(), "failed to connect to execution client")

	// Verify no metrics were registered
	metrics, err := reg.Gather()
	require.NoError(t, err)
	assert.Empty(t, metrics)
}

func TestIntegration_DiskUsage(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		DiskUsage: DiskUsageConfig{
			Enabled:     true,
			Directories: []string{"/tmp"},
			Interval:    "1m",
		},
	}
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()

	integration := New(logger, cfg, reg)

	// Start the integration
	err := integration.Run(context.Background())
	require.NoError(t, err)

	// Give some time for metrics to be collected
	time.Sleep(100 * time.Millisecond)

	// Verify metrics were registered
	metrics, err := reg.Gather()
	require.NoError(t, err)
	assert.NotEmpty(t, metrics)

	// Stop the integration
	integration.Stop()
}

func TestIntegration_InvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
	}{
		{
			name: "invalid disk usage interval",
			cfg: &Config{
				Enabled: true,
				DiskUsage: DiskUsageConfig{
					Enabled:     true,
					Directories: []string{"/tmp"},
					Interval:    "invalid",
				},
			},
		},
		{
			name: "empty execution URL",
			cfg: &Config{
				Enabled: true,
				Execution: ExecutionConfig{
					Enabled: true,
					URL:     "",
				},
			},
		},
		{
			name: "empty consensus URL",
			cfg: &Config{
				Enabled: true,
				Consensus: ConsensusConfig{
					Enabled: true,
					URL:     "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			reg := prometheus.NewRegistry()

			integration := New(logger, tt.cfg, reg)
			err := integration.Run(context.Background())
			assert.Error(t, err)
		})
	}
}

func TestIntegration_Stop(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		DiskUsage: DiskUsageConfig{
			Enabled:     true,
			Directories: []string{"/tmp"},
			Interval:    "1m",
		},
	}
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()

	integration := New(logger, cfg, reg)

	// Start the integration
	err := integration.Run(context.Background())
	require.NoError(t, err)

	// Stop the integration
	integration.Stop()

	// Verify metrics are no longer being collected
	metrics, err := reg.Gather()
	require.NoError(t, err)
	assert.Empty(t, metrics)
}
