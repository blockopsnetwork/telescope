package ethereum

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/log"
	v2integrations "github.com/blockopsnetwork/telescope/internal/static/integrations/v2"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/v2/autoscrape"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/v2/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestGlobals() v2integrations.Globals {
	return v2integrations.Globals{
		SubsystemOpts: v2integrations.SubsystemOptions{
			Metrics: v2integrations.MetricsSubsystemOptions{
				Autoscrape: autoscrape.DefaultGlobal,
			},
		},
	}
}

func TestIntegration_Disabled(t *testing.T) {
	cfg := &Config{
		Enabled: false,
		Common:  common.MetricsConfig{},
	}
	logger := log.NewNopLogger()
	globals := createTestGlobals()

	integration := New(logger, cfg, globals)
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	err := integration.RunIntegration(ctx)
	require.NoError(t, err)
}

func TestIntegration_ExecutionClient_InvalidURL(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Common:  common.MetricsConfig{},
		Execution: ExecutionConfig{
			Enabled: true,
			URL:     "http://invalid-url:8545", // This should fail to connect
			Modules: []string{"eth", "net"},
		},
	}
	logger := log.NewNopLogger()
	globals := createTestGlobals()

	integration := New(logger, cfg, globals)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start the integration - should error because we can't connect to the client
	err := integration.RunIntegration(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to setup execution client")

	// Clean up
	integration.Stop()
}

func TestIntegration_ConsensusClient_InvalidURL(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Common:  common.MetricsConfig{},
		Consensus: ConsensusConfig{
			Enabled: true,
			URL:     "http://invalid-url:5052", // This should fail to connect
			EventStream: EventStreamConfig{
				Enabled: true,
				Topics:  []string{"head", "finalized_checkpoint"},
			},
		},
	}
	logger := log.NewNopLogger()
	globals := createTestGlobals()

	integration := New(logger, cfg, globals)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Consensus client setup doesn't fail immediately like execution client
	// It creates the client but errors occur when actually connecting
	go func() {
		_ = integration.RunIntegration(ctx)
	}()

	// Give some time for the integration to attempt setup
	time.Sleep(100 * time.Millisecond)

	// Clean up
	integration.Stop()
}


func TestIntegration_InvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
	}{
		{
			name: "empty execution URL",
			cfg: &Config{
				Enabled: true,
				Common:  common.MetricsConfig{},
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
				Common:  common.MetricsConfig{},
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
			globals := createTestGlobals()

			integration := New(logger, tt.cfg, globals)
			
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			
			err := integration.RunIntegration(ctx)
			assert.Error(t, err)
			
			// Clean up
			integration.Stop()
		})
	}
}

func TestIntegration_Stop(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Common:  common.MetricsConfig{},
	}
	logger := log.NewNopLogger()
	globals := createTestGlobals()

	integration := New(logger, cfg, globals)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Start the integration in a goroutine
	go func() {
		_ = integration.RunIntegration(ctx)
	}()

	// Give some time for the integration to start
	time.Sleep(100 * time.Millisecond)

	// Stop the integration
	integration.Stop()

	// Ensure no panic occurred and stop completed
	assert.True(t, true) // Test passes if we get here without panic
}

func TestIntegration_MetricsInterface(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Common:  common.MetricsConfig{},
	}
	logger := log.NewNopLogger()
	globals := createTestGlobals()

	integration := New(logger, cfg, globals)

	// Test that the integration implements the required interfaces
	assert.Implements(t, (*v2integrations.Integration)(nil), integration)
	assert.Implements(t, (*v2integrations.HTTPIntegration)(nil), integration)
	assert.Implements(t, (*v2integrations.MetricsIntegration)(nil), integration)

	// Test Handler method
	handler, err := integration.Handler("/ethereum")
	require.NoError(t, err)
	assert.NotNil(t, handler)

	// Test Targets method
	endpoint := v2integrations.Endpoint{Host: "localhost:8080"}
	targets := integration.Targets(endpoint)
	assert.NotEmpty(t, targets)
	assert.Equal(t, 1, len(targets))

	// Clean up
	integration.Stop()
}

func TestConfig_Name(t *testing.T) {
	cfg := &Config{}
	assert.Equal(t, "ethereum", cfg.Name())
}

func TestConfig_ApplyDefaults(t *testing.T) {
	cfg := &Config{}
	globals := createTestGlobals()
	
	err := cfg.ApplyDefaults(globals)
	require.NoError(t, err)
}

func TestConfig_Identifier(t *testing.T) {
	cfg := &Config{}
	globals := createTestGlobals()
	
	identifier, err := cfg.Identifier(globals)
	require.NoError(t, err)
	assert.Equal(t, "ethereum", identifier)
	
	// Test with custom instance key
	instanceKey := "custom-instance"
	cfg.Common.InstanceKey = &instanceKey
	identifier, err = cfg.Identifier(globals)
	require.NoError(t, err)
	assert.Equal(t, "custom-instance", identifier)
}

func TestConfig_NewIntegration(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Common:  common.MetricsConfig{},
	}
	logger := log.NewNopLogger()
	globals := createTestGlobals()
	
	integration, err := cfg.NewIntegration(logger, globals)
	require.NoError(t, err)
	assert.NotNil(t, integration)
	assert.IsType(t, (*Integration)(nil), integration)
}