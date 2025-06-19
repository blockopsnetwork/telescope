package main

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	networksConfig "github.com/blockopsnetwork/telescope/internal/static/config/networks"

	"github.com/go-kit/log/level"
	"github.com/blockopsnetwork/telescope/internal/boringcrypto"
	"github.com/blockopsnetwork/telescope/internal/build"
	"github.com/blockopsnetwork/telescope/internal/static/config"
	"github.com/blockopsnetwork/telescope/internal/static/server"
	util_log "github.com/blockopsnetwork/telescope/internal/util/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/grafana/loki/clients/pkg/promtail/discovery/consulagent"
	_ "github.com/prometheus/prometheus/discovery/install"
	_ "github.com/blockopsnetwork/telescope/internal/static/integrations/cadvisor"

	_ "github.com/blockopsnetwork/telescope/internal/static/integrations/install"
	_ "golang.org/x/crypto/x509roots/fallback"
)

var cmd = &cobra.Command{
	Use:   "telescope",
	Short: "An All-in-One Web3 Observability tooling",
	Long:  `Gain full insights into the performance of your dApps, nodes and onchain events with Telescope.`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultCfg := server.DefaultConfig()
		logger := server.NewLogger(&defaultCfg)

		// Check if config file is specified
		configFile := viper.GetString("config-file")
		if configFile != "" {
			if err := validateConfigFile(configFile); err != nil {
				level.Error(logger).Log("msg", "invalid config file", "err", err)
				os.Exit(1)
			}
			level.Info(logger).Log("msg", "using config file", "path", configFile)
			agent(configFile)
			return
		}

		// Otherwise load config from flags/env
		var config TelescopeConfig
		if err := config.loadConfig(); err != nil {
			level.Error(logger).Log("msg", "failed to load config", "err", err)
			os.Exit(1)
		}

		// Get network config and generate scrape configs
		networkConfig := getNetworkConfig(config.Network)
		level.Info(logger).Log("msg", "starting telescope agent", "network", config.Network)
		scrapeConfigs := networkConfig.GenerateScrapeConfigs(config.ProjectName, config.Network)

		// Generate and write full config
		fullConfig := generateFullConfig(config, scrapeConfigs)
		configFilePath := "telescope_config.yaml"
		if err := writeConfigToFile(fullConfig, configFilePath); err != nil {
			level.Error(logger).Log("msg", "failed to write config file", "path", configFilePath, "err", err)
			os.Exit(1)
		}

		level.Info(logger).Log("msg", "configuration written", "path", configFilePath)
		agent(configFilePath)
	},
}

// getSupportedNetworks returns a sorted list of supported blockchain networks.
func getSupportedNetworks() []string {
	networks := make([]string, 0, len(networkConfigs))
	for network := range networkConfigs {
		networks = append(networks, network)
	}
	sort.Strings(networks)
	return networks
}


var networkConfigs = map[string]networksConfig.NetworkConfig{
	"ethereum":    networksConfig.NewEthereumConfig(),
	"polkadot":    networksConfig.NewPolkadotConfig(),
	"hyperbridge": networksConfig.NewHyperbridgeConfig(),
	"ssv":         networksConfig.NewSSVConfig(),
}

type Config struct {
	Server       ServerConfig           `yaml:"server"`
	Metrics      MetricsConfig          `yaml:"metrics"`
	Logs         LogsConfig             `yaml:"logs"`
	Integrations map[string]interface{} `yaml:"integrations"`
}

type ServerConfig struct {
	LogLevel string `yaml:"log_level"`
}

type MetricsConfig struct {
	Global        GlobalConfig   `yaml:"global"`
	Wal_Directory string         `yaml:"wal_directory"`
	Configs       []MetricConfig `yaml:"configs"`
}

type GlobalConfig struct {
	ScrapeInterval string            `yaml:"scrape_interval"`
	ExternalLabels map[string]string `yaml:"external_labels"`
	RemoteWrite    []RemoteWrite     `yaml:"remote_write"`
}

type MetricConfig struct {
	Name          string         `yaml:"name"`
	HostFilter    bool           `yaml:"host_filter"`
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs"`
}

type ScrapeConfig struct {
	JobName       string         `yaml:"job_name"`
	StaticConfigs []StaticConfig `yaml:"static_configs"`
}

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type IntegrationsConfig struct {
    Agent        ToIntegrate         `yaml:"agent"`
    NodeExporter ToIntegrate         `yaml:"node_exporter"`
    Cadvisor     CadvisorIntegration `yaml:"cadvisor"`
}

// Add a new struct for cAdvisor specific configuration
type CadvisorIntegration struct {
    Enabled    bool `yaml:"enabled"`
    DockerOnly bool `yaml:"docker_only"`
}

type ToIntegrate struct {
	Enabled bool `yaml:"enabled"`
}

type StaticConfig struct {
	Targets []string `yaml:"targets"`
}

type RemoteWrite struct {
	URL       string    `yaml:"url"`
	BasicAuth BasicAuth `yaml:"basic_auth"`
}

type LogsConfig struct {
	Configs []LogConfig `yaml:"configs"`
}

type LogConfig struct {
	Name      string      `yaml:"name"`
	Clients   []LogClient `yaml:"clients"`
	Positions Positions   `yaml:"positions"`
}

type Positions struct {
	Filename string `yaml:"filename"`
}

type LogClient struct {
	URL            string            `yaml:"url"`
	BasicAuth      BasicAuth         `yaml:"basic_auth"`
	ExternalLabels map[string]string `yaml:"external_labels"`
}

type TelescopeConfig struct {
	Metrics           bool
	Logs              bool
	Network           string
	ProjectId         string
	ProjectName       string
	TelescopeUsername string
	TelescopePassword string
	RemoteWriteUrl    string
	LokiUsername      string
	LokiPassword      string
	LogsSinkURL       string
	// Ethereum integration fields
	EthereumEnabled            bool
	EthereumExecutionURL       string
	EthereumConsensusURL       string
	EthereumExecutionModules   []string
	EthereumDiskUsageEnabled   bool
	EthereumDiskUsageDirs      []string
	EthereumDiskUsageInterval  string
}

func handleErr(err error, msg string) {
	if err != nil {
		defaultCfg := server.DefaultConfig()
		logger := server.NewLogger(&defaultCfg)
		level.Error(logger).Log("msg", msg, "err", err)
		os.Exit(1)
	}
}

func toLowerAndEscape(input string) string {
	lowercase := strings.ToLower(input)
	escaped := url.QueryEscape(lowercase)
	return escaped
}

func getNetworkConfig(network string) networksConfig.NetworkConfig {
	config, exists := networkConfigs[network]
	if !exists {
		supportedNetworks := getSupportedNetworks()
		handleErr(
			fmt.Errorf("unsupported network %q, must be one of: %s",
				network,
				strings.Join(supportedNetworks, ", ")),
			"Invalid network configuration",
		)
	}
	return config
}

// performs basic validation of the TelescopeConfig.
// It checks URLs, required fields, and validates the network selection.
func (c *TelescopeConfig) validate() error {
	// Validate URLs
	if c.RemoteWriteUrl != "" {
		if err := validateURL(c.RemoteWriteUrl, "remote write"); err != nil {
			return err
		}
	}

	if c.LogsSinkURL != "" {
		if err := validateURL(c.LogsSinkURL, "logs sink"); err != nil {
			return err
		}
	}

	// Validate other required fields
	if c.ProjectId == "" {
		return fmt.Errorf("project ID cannot be empty")
	}

	if c.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Validate Network is supported
	if _, exists := networkConfigs[c.Network]; !exists {
		supportedNetworks := getSupportedNetworks()
		return fmt.Errorf("unsupported network %q, must be one of: %s",
			c.Network,
			strings.Join(supportedNetworks, ", "))
	}

	return nil
}

// checks if a URL string is valid and properly formatted.
// It verifies that the URL has both a scheme (http/https) and a host.
func validateURL(urlStr, name string) error {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid %s URL: %v", name, err)
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("%s URL must include scheme (http:// or https://)", name)
	}
	if parsed.Host == "" {
		return fmt.Errorf("%s URL must include host", name)
	}
	return nil
}

// checks if the provided config file path is valid,
// exists, and points to a regular file rather than a directory.
func validateConfigFile(configFile string) error {
	if configFile == "" {
		return fmt.Errorf("config file path cannot be empty")
	}

	info, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist: %s", configFile)
		}
		return fmt.Errorf("error accessing config file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("config file path points to a directory: %s", configFile)
	}

	return nil
}
// validates the logs-specific configuration.
// It checks that all required fields for log collection are provided when logs are enabled.
func (c *TelescopeConfig) validateLogsConfig() error {
	if !c.Logs {
		return nil // Logs not enabled, no validation needed
	}

	if c.LogsSinkURL == "" {
		return fmt.Errorf("logs-sink-url is required when logs are enabled")
	}

	if c.LokiUsername == "" {
		return fmt.Errorf("telescope-loki-username is required when logs are enabled")
	}

	if c.LokiPassword == "" {
		return fmt.Errorf("telescope-loki-password is required when logs are enabled")
	}

	return nil
}

// validates the metrics-specific configuration.
// It checks that all required fields for metrics collection are provided when metrics are enabled.
func (c *TelescopeConfig) validateMetricsConfig() error {
	if !c.Metrics {
		return nil // Metrics not enabled, no validation needed
	}

	if c.TelescopeUsername == "" {
		return fmt.Errorf("telescope-username is required when metrics are enabled")
	}

	if c.TelescopePassword == "" {
		return fmt.Errorf("telescope-password is required when metrics are enabled")
	}

	if c.RemoteWriteUrl == "" {
		return fmt.Errorf("remote-write-url is required when metrics are enabled")
	}

	return nil
}

// generates the complete agent configuration from the provided
// TelescopeConfig and network scrape configurations.
func generateFullConfig(config TelescopeConfig, networkScrapeConfigs []networksConfig.ScrapeConfig) Config {
	scrapeConfigs := make([]ScrapeConfig, len(networkScrapeConfigs))
	for i, nsc := range networkScrapeConfigs {
		scrapeConfigs[i] = ScrapeConfig{
			JobName: nsc.JobName,
			StaticConfigs: []StaticConfig{
				{
					Targets: nsc.StaticConfigs[0].Targets,
				},
			},
		}
	}

	metricsInstanceName := toLowerAndEscape(config.ProjectName + "_" + config.Network + "_metrics")
	integrations := map[string]interface{}{
		"agent": map[string]interface{}{
			"autoscrape": map[string]interface{}{
				"enable":           true,
				"metrics_instance": metricsInstanceName,
			},
		},
		"node_exporter": map[string]interface{}{
			"autoscrape": map[string]interface{}{
				"enable":           true,
				"metrics_instance": metricsInstanceName,
			},
		},
	}

	// Add Ethereum integration if enabled
	if config.EthereumEnabled || config.EthereumExecutionURL != "" || config.EthereumConsensusURL != "" || config.EthereumDiskUsageEnabled {

		ethereumConfig := map[string]interface{}{
			"instance": "ethereum_node_1",
			"enabled":  true,
			"autoscrape": map[string]interface{}{
				"enable":           true,
				"metrics_instance": metricsInstanceName,
			},
		}

		// Add execution config if URL provided
		if config.EthereumExecutionURL != "" {
			modules := config.EthereumExecutionModules
			if len(modules) == 0 {
				modules = []string{"sync", "eth", "net", "web3", "txpool"}
			}
			ethereumConfig["execution"] = map[string]interface{}{
				"enabled": true,
				"url":     config.EthereumExecutionURL,
				"modules": modules,
			}
		}

		// Add consensus config if URL provided
		if config.EthereumConsensusURL != "" {
			ethereumConfig["consensus"] = map[string]interface{}{
				"enabled": true,
				"url":     config.EthereumConsensusURL,
				"event_stream": map[string]interface{}{
					"enabled": true,
					"topics":  []string{"head", "finalized_checkpoint"},
				},
			}
		}

		// Add disk usage config if enabled
		if config.EthereumDiskUsageEnabled {
			interval := config.EthereumDiskUsageInterval
			if interval == "" {
				interval = "5m" // default
			}
			ethereumConfig["disk_usage"] = map[string]interface{}{
				"enabled":     true,
				"directories": config.EthereumDiskUsageDirs,
				"interval":    interval,
			}
		}

		integrations["ethereum_configs"] = []interface{}{ethereumConfig}
	}

	cfg := Config{
		Server: ServerConfig{
			LogLevel: "info",
		},
		Metrics: MetricsConfig{
			Wal_Directory: "/tmp/telescope",
			Global: GlobalConfig{
				ScrapeInterval: "15s",
				ExternalLabels: map[string]string{
					"project_id":   config.ProjectId,
					"project_name": config.ProjectName,
				},
				RemoteWrite: []RemoteWrite{
					{
						URL: config.RemoteWriteUrl,
						BasicAuth: BasicAuth{
							Username: config.TelescopeUsername,
							Password: config.TelescopePassword,
						},
					},
				},
			},
			Configs: []MetricConfig{
				{
					Name:          toLowerAndEscape(config.ProjectName + "_" + config.Network + "_metrics"),
					HostFilter:    false,
					ScrapeConfigs: scrapeConfigs,
				},
			},
		},
		Integrations: integrations,
	}

	if config.Logs {
		cfg.Logs = LogsConfig{
			Configs: []LogConfig{
				{
					Name: "telescope_logs",
					Clients: []LogClient{
						{
							URL: config.LogsSinkURL,
							BasicAuth: BasicAuth{
								Username: config.LokiUsername,
								Password: config.LokiPassword,
							},
							ExternalLabels: map[string]string{
								"project_id":   config.ProjectId,
								"project_name": config.ProjectName,
							},
						},
					},
					Positions: Positions{
						Filename: "/tmp/telescope_logs",
					},
				},
			},
		}
	}

	return cfg
}

// writes the provided Config to a YAML file at the specified path.
func writeConfigToFile(config Config, filePath string) error {
	if filePath == "" {
		return fmt.Errorf("empty file path provided")
	}

	var buf strings.Builder
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(&config); err != nil {
		return fmt.Errorf("error marshaling to YAML: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

// loads and validates the agent configuration from command line flags.
// It checks for required flags based on enabled features and validates the configuration.
func (c *TelescopeConfig) loadConfig() error {
	if err := checkRequiredFlags(); err != nil {
		return fmt.Errorf("flag validation error: %w", err)
	}

	// Load values from viper
	c.Metrics = viper.GetBool("metrics")
	c.Logs = viper.GetBool("enable-logs")
	c.Network = viper.GetString("network")
	c.ProjectId = viper.GetString("project-id")
	c.ProjectName = viper.GetString("project-name")
	c.TelescopeUsername = viper.GetString("telescope-username")
	c.TelescopePassword = viper.GetString("telescope-password")
	c.RemoteWriteUrl = viper.GetString("remote-write-url")
	c.LokiUsername = viper.GetString("telescope-loki-username")
	c.LokiPassword = viper.GetString("telescope-loki-password")
	c.LogsSinkURL = viper.GetString("logs-sink-url")
	
	// Load Ethereum integration values
	c.EthereumEnabled = viper.GetBool("ethereum-enabled")
	c.EthereumExecutionURL = viper.GetString("ethereum-execution-url")
	c.EthereumConsensusURL = viper.GetString("ethereum-consensus-url")
	c.EthereumExecutionModules = viper.GetStringSlice("ethereum-execution-modules")
	c.EthereumDiskUsageEnabled = viper.GetBool("ethereum-disk-usage-enabled")
	c.EthereumDiskUsageDirs = viper.GetStringSlice("ethereum-disk-usage-dirs")
	c.EthereumDiskUsageInterval = viper.GetString("ethereum-disk-usage-interval")

	// Run all validations
	if err := c.validate(); err != nil {
		return fmt.Errorf("config validation error: %w", err)
	}

	if err := c.validateMetricsConfig(); err != nil {
		return fmt.Errorf("metrics configuration error: %w", err)
	}

	if err := c.validateLogsConfig(); err != nil {
		return fmt.Errorf("logs configuration error: %w", err)
	}

	if err := c.validateEthereumConfig(); err != nil {
		return fmt.Errorf("ethereum configuration error: %w", err)
	}

	return nil
}

// validates ethereum integration configuration.
// Ensures that if ethereum flags are provided, integrations-next feature is enabled.
func (c *TelescopeConfig) validateEthereumConfig() error {
	// Check if any ethereum flags are provided
	ethereumEnabled := c.EthereumEnabled || c.EthereumExecutionURL != "" || c.EthereumConsensusURL != "" || c.EthereumDiskUsageEnabled
	
	if ethereumEnabled {
		enableFeatures := viper.GetString("enable-features")
		if !strings.Contains(enableFeatures, "integrations-next") {
			return fmt.Errorf("ethereum integration requires --enable-features integrations-next flag to be specified")
		}
		
		// Validate that at least one ethereum component is configured
		if c.EthereumExecutionURL == "" && c.EthereumConsensusURL == "" && !c.EthereumDiskUsageEnabled {
			return fmt.Errorf("when ethereum integration is enabled, at least one of --ethereum-execution-url, --ethereum-consensus-url, or --ethereum-disk-usage-enabled must be provided")
		}
		
		// Validate URLs if provided
		if c.EthereumExecutionURL != "" {
			if err := validateURL(c.EthereumExecutionURL, "ethereum execution"); err != nil {
				return err
			}
		}
		if c.EthereumConsensusURL != "" {
			if err := validateURL(c.EthereumConsensusURL, "ethereum consensus"); err != nil {
				return err
			}
		}
		
		// Validate disk usage configuration
		if c.EthereumDiskUsageEnabled {
			if len(c.EthereumDiskUsageDirs) == 0 {
				return fmt.Errorf("when --ethereum-disk-usage-enabled is true, --ethereum-disk-usage-dirs must be provided")
			}
			// Validate interval format
			if c.EthereumDiskUsageInterval != "" {
				if _, err := time.ParseDuration(c.EthereumDiskUsageInterval); err != nil {
					return fmt.Errorf("invalid --ethereum-disk-usage-interval format: %w", err)
				}
			}
		}
	}
	
	return nil
}

// verifies that all required command line flags are provided.
// Required flags depend on which features (metrics/logs) are enabled.
func checkRequiredFlags() error {
	// Always required
	baseFlags := []string{
		"network",
		"project-id",
		"project-name",
	}

	// Required for metrics
	metricsFlags := []string{
		"telescope-username",
		"telescope-password",
		"remote-write-url",
	}

	// Required for logs
	logsFlags := []string{
		"logs-sink-url",
		"telescope-loki-username",
		"telescope-loki-password",
	}

	missingFlags := []string{}

	// Check base flags
	for _, flag := range baseFlags {
		if viper.GetString(flag) == "" {
			missingFlags = append(missingFlags, flag)
		}
	}

	// Check metrics flags if metrics enabled
	if viper.GetBool("metrics") {
		for _, flag := range metricsFlags {
			if viper.GetString(flag) == "" {
				missingFlags = append(missingFlags, flag)
			}
		}
	}

	// Check logs flags if logs enabled
	if viper.GetBool("enable-logs") {
		for _, flag := range logsFlags {
			if viper.GetString(flag) == "" {
				missingFlags = append(missingFlags, flag)
			}
		}
	}

	if len(missingFlags) > 0 && viper.GetString("config-file") == "" {
		return fmt.Errorf("missing required flags: %s", strings.Join(missingFlags, ", "))
	}

	return nil
}

func init() {
	prometheus.MustRegister(build.NewCollector("agent"))
	cobra.OnInitialize(initConfig)

	// Set command description and examples
	cmd.Use = "telescope"
	cmd.Short = "Telescope monitoring agent"
	cmd.Long = `Telescope is a monitoring agent that collects metrics and logs from various sources
and forwards them to remote endpoints. It supports multiple integrations including
Ethereum blockchain monitoring, system metrics, and custom applications.

IMPORTANT: Ethereum integration requires --enable-features integrations-next`

	cmd.Example = `  # Basic metrics collection
  telescope --network=ethereum --project-id=my-project --project-name=my-project \
            --telescope-username=user --telescope-password=pass \
            --remote-write-url=https://prometheus.example.com/api/v1/write

  # With Ethereum integration  
  telescope --enable-features integrations-next --network=ethereum \
            --project-id=my-project --project-name=my-project \
            --telescope-username=user --telescope-password=pass \
            --remote-write-url=https://prometheus.example.com/api/v1/write \
            --ethereum-execution-url=http://localhost:8545 \
            --ethereum-consensus-url=http://localhost:5052

  # With disk usage monitoring (separate from node_exporter)
  telescope --enable-features integrations-next --network=ssv \
            --project-id=my-project --project-name=my-project \
            --telescope-username=user --telescope-password=pass \
            --remote-write-url=https://prometheus.example.com/api/v1/write \
            --ethereum-execution-url=http://localhost:8545 \
            --ethereum-disk-usage-enabled \
            --ethereum-disk-usage-dirs=/data/ethereum,/data/consensus \
            --ethereum-disk-usage-interval=10m`

	// Basic configuration flags
	cmd.Flags().String("config-file", "", "Config file path (alternative to using flags)")
	cmd.Flags().String("network", "", fmt.Sprintf("Target network (%s)", strings.Join(getSupportedNetworks(), ", ")))
	cmd.Flags().String("project-id", "", "Project identifier")
	cmd.Flags().String("project-name", "", "Project name for labeling")

	// Metrics configuration flags  
	cmd.Flags().Bool("metrics", true, "Enable metrics collection")
	cmd.Flags().String("telescope-username", "", "Username for remote write authentication")
	cmd.Flags().String("telescope-password", "", "Password for remote write authentication")
	cmd.Flags().String("remote-write-url", "", "Prometheus remote write endpoint URL")

	// Logs configuration flags
	cmd.Flags().Bool("enable-logs", false, "Enable log collection")
	cmd.Flags().String("logs-sink-url", "", "Log sink endpoint URL")
	cmd.Flags().String("telescope-loki-username", "", "Username for Loki authentication")
	cmd.Flags().String("telescope-loki-password", "", "Password for Loki authentication")

	// Feature flags
	cmd.Flags().String("enable-features", "", "Experimental features (comma-separated, e.g., integrations-next)")
	
	// Ethereum integration flags
	cmd.Flags().Bool("ethereum-enabled", false, "Enable Ethereum metrics collection")
	cmd.Flags().String("ethereum-execution-url", "", "Ethereum execution node URL (e.g., http://localhost:8545)")
	cmd.Flags().String("ethereum-consensus-url", "", "Ethereum consensus node URL (e.g., http://localhost:5052)")
	cmd.Flags().StringSlice("ethereum-execution-modules", []string{"sync", "eth", "net", "web3", "txpool"}, "Execution modules to enable (comma-separated)")
	cmd.Flags().Bool("ethereum-disk-usage-enabled", false, "Enable Ethereum disk usage monitoring")
	cmd.Flags().StringSlice("ethereum-disk-usage-dirs", []string{}, "Directories to monitor for Ethereum disk usage (comma-separated)")
	cmd.Flags().String("ethereum-disk-usage-interval", "5m", "Interval for disk usage collection (e.g., 1h, 5m, 30s)")

	// Mark required flags
	cmd.MarkFlagRequired("network")
	cmd.MarkFlagRequired("project-id") 
	cmd.MarkFlagRequired("project-name")

	// Bind all flags to viper
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

// initializes and runs the monitoring agent with the provided configuration file.
func agent(configPath string) {
	defaultCfg := server.DefaultConfig()
	logger := server.NewLogger(&defaultCfg)

	reloader := func(log *server.Logger) (*config.Config, error) {
		fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		// Add enable-features flag support for v2 integrations
		args := []string{"-config.file", configPath}
		enableFeatures := viper.GetString("enable-features")
		if enableFeatures != "" {
			args = append(args, "-enable-features", enableFeatures)
		}
		return config.Load(fs, args, log)
	}

	cfg, err := reloader(logger)
	if err != nil {
		level.Error(logger).Log("msg", "failed to load config", "err", err)
		os.Exit(1)
	}

	logger = server.NewLogger(cfg.Server)
	util_log.Logger = logger

	level.Info(logger).Log("msg", "starting agent", "boringcrypto", boringcrypto.Enabled)
	ep, err := NewEntrypoint(logger, cfg, reloader)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create agent entrypoint", "err", err)
		os.Exit(1)
	}

	if err = ep.Start(); err != nil {
		level.Error(logger).Log("msg", "error running agent", "err", err)
		ep.Stop() // Ensure cleanup happens
		os.Exit(1)
	}

	ep.Stop()
	level.Info(logger).Log("msg", "agent exiting")
}

func main() {
	cobra.OnInitialize(initConfig)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
