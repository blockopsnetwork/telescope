package ethereum

// Config holds the configuration for the ethereum integration
type Config struct {
	Enabled bool `yaml:"enabled"`

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
	Enabled  bool   `yaml:"enabled"`
	URL      string `yaml:"url"`
	Timeout  string `yaml:"timeout"`
	Interval string `yaml:"interval"`
}

// DiskUsageConfig holds the configuration for disk usage monitoring
type DiskUsageConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Directories []string `yaml:"directories"`
	Interval    string   `yaml:"interval"`
}

// DefaultConfig returns the default configuration for the ethereum integration
func DefaultConfig() Config {
	return Config{
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
		},
		DiskUsage: DiskUsageConfig{
			Enabled:     false,
			Directories: []string{"/var/lib/ethereum"},
			Interval:    "5m",
		},
	}
}

// UnmarshalRiver implements river.Unmarshaler
func (c *Config) UnmarshalRiver(f func(interface{}) error) error {
	*c = DefaultConfig()

	type config Config
	return f((*config)(c))
}
