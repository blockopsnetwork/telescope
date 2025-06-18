package ethereum

import (
	"github.com/go-kit/log"
	v2 "github.com/blockopsnetwork/telescope/internal/static/integrations/v2"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/v2/common"
)

// DefaultConfig is the default configuration for the ethereum integration
var DefaultConfig = Config{
	Enabled: false,
	Execution: ExecutionConfig{
		Enabled:  false,
		URL:      "http://localhost:8545",
		Modules:  []string{"eth", "net", "web3"},
		Timeout:  "5s",
		Interval: "15s",
	},
	Consensus: ConsensusConfig{
		Enabled:  false,
		URL:      "http://localhost:5052",
		Timeout:  "5s",
		Interval: "15s",
		EventStream: EventStreamConfig{
			Enabled: false,
			Topics:  []string{},
		},
	},
	DiskUsage: DiskUsageConfig{
		Enabled:     false,
		Directories: []string{"/var/lib/ethereum"},
		Interval:    "5m",
	},
}

// Config holds the configuration for the ethereum integration
type Config struct {
	Common  common.MetricsConfig `yaml:",inline"`
	Enabled bool                 `yaml:"enabled"`

	// Execution client configuration
	Execution ExecutionConfig `yaml:"execution"`

	// Consensus client configuration
	Consensus ConsensusConfig `yaml:"consensus"`

	// Disk usage monitoring configuration
	DiskUsage DiskUsageConfig `yaml:"disk_usage"`
}

// ExecutionConfig holds the configuration for the execution client
type ExecutionConfig struct {
	Enabled  bool     `yaml:"enabled"`
	URL      string   `yaml:"url"`
	Modules  []string `yaml:"modules"`
	Timeout  string   `yaml:"timeout"`
	Interval string   `yaml:"interval"`
}

// ConsensusConfig holds the configuration for the consensus client
type ConsensusConfig struct {
	Enabled     bool              `yaml:"enabled"`
	URL         string            `yaml:"url"`
	Timeout     string            `yaml:"timeout"`
	Interval    string            `yaml:"interval"`
	EventStream EventStreamConfig `yaml:"event_stream"`
}

// EventStreamConfig holds the configuration for the event stream
type EventStreamConfig struct {
	Enabled bool     `yaml:"enabled"`
	Topics  []string `yaml:"topics"`
}

// DiskUsageConfig holds the configuration for disk usage monitoring
type DiskUsageConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Directories []string `yaml:"directories"`
	Interval    string   `yaml:"interval"`
}

// Name returns the name of the integration
func (c *Config) Name() string {
	return "ethereum"
}

// ApplyDefaults applies default values to the Config
func (c *Config) ApplyDefaults(g v2.Globals) error {
	c.Common.ApplyDefaults(g.SubsystemOpts.Metrics.Autoscrape)
	return nil
}

// Identifier returns a unique identifier for the integration
func (c *Config) Identifier(g v2.Globals) (string, error) {
	// Use a combination of integration name and instance name
	if c.Common.InstanceKey != nil {
		return *c.Common.InstanceKey, nil
	}
	return c.Name(), nil
}

// NewIntegration creates a new integration from the config
func (c *Config) NewIntegration(l log.Logger, g v2.Globals) (v2.Integration, error) {
	return New(l, c, g), nil
}

// UnmarshalYAML implements yaml.Unmarshaler for Config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig

	type config Config
	return unmarshal((*config)(c))
}

