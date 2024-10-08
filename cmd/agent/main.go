package main

import (

    "bytes"
    "flag"
    "fmt"
    "log"
    "net/url"
    "os"
    "strings"

    networksConfig "github.com/grafana/agent/internal/static/config/networks"

    "github.com/go-kit/log/level"
    "github.com/grafana/agent/internal/boringcrypto"
    "github.com/grafana/agent/internal/build"
    "github.com/grafana/agent/internal/static/config"
    "github.com/grafana/agent/internal/static/server"
    util_log "github.com/grafana/agent/internal/util/log"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "gopkg.in/yaml.v3"

    "github.com/prometheus/client_golang/prometheus"

    // Register Prometheus SD components
    _ "github.com/grafana/loki/clients/pkg/promtail/discovery/consulagent"
    _ "github.com/prometheus/prometheus/discovery/install"

    // Register integrations
    _ "github.com/grafana/agent/internal/static/integrations/install"

    // Embed a set of fallback X.509 trusted roots
    // Allows the app to work correctly even when the OS does not provide a verifier or systems roots pool
    _ "golang.org/x/crypto/x509roots/fallback"
)

var network string

var cmd = &cobra.Command{
    Use:   "telescope",
    Short: "An All-in-One Web3 Observability tooling",
    Long:  `Gain full insights into the performance of your dApps, nodes and onchain events with Telescope.`,
    Run: func(cmd *cobra.Command, args []string) {
        var config TelescopeConfig

        if err := config.loadConfig(); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // Proceed to generate and write configuration, then start the agent
        network := viper.GetString("network")
        projectName := viper.GetString("project-name")
        networkConfig := getNetworkConfig(network)

        fmt.Println("Starting Telescope Agent with network config:", networkConfig)
        scrapeConfigs := networkConfig.GenerateScrapeConfigs(projectName, network)
        fullConfig := generateFullConfig(config, scrapeConfigs)
        configFilePath := "telescope_config.yaml"
        if err := writeConfigToFile(fullConfig, configFilePath); err != nil {
            log.Fatalf("Failed to write config to file: %v", err)
        }
		
        fmt.Printf("Configuration written to %s\n", configFilePath)

        agent(configFilePath)
    },
}

var helperFunction = func(cmd *cobra.Command, args []string) {
    fmt.Println("Custom Help:")
    fmt.Println("Usage: telescope start [flags]")
    fmt.Println("\nFlags:")
    fmt.Println("  --metrics\t\tEnable metrics")
    fmt.Println("  --network\t\tSpecify the network")
    fmt.Println("  --project-id\t\tSpecify the project ID")
    fmt.Println("  --project-name\tSpecify the project name")
    fmt.Println("  --telescope-username\tSpecify the telescope username")
    fmt.Println("  --telescope-password\tSpecify the telescope password")
    fmt.Println("  --remote-write-url\tSpecify the remote write URL")
    fmt.Println("  --config-file\t\tSpecify the config file")
    fmt.Println("  --enable-logs\t\tEnable logs")
    fmt.Println("  --logs-sink-url\tSpecify the Log Sink URL")
    fmt.Println("  --telescope-loki-username\tSpecify the Loki username")
    fmt.Println("  --telescope-loki-password\tSpecify the Loki password")
}

var cfgFile string

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
    Agent        ToIntegrate `yaml:"agent"`
    NodeExporter ToIntegrate `yaml:"node_exporter"`
}

type ToIntegrate struct {
    Enabled bool `yaml:"enabled"`
}

type StaticConfig struct {
    Targets []string `yaml:"targets"`
}

type RemoteWrite struct {
    URL       string            `yaml:"url"`
    BasicAuth map[string]string `yaml:"basic_auth"`
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
}

func handleErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %v", msg, err)
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
        handleErr(fmt.Errorf("unsupported network: %s", network), "Invalid network")
    }
    return config
}

func generateFullConfig(config TelescopeConfig, networkScrapeConfigs []networksConfig.ScrapeConfig) Config {
    cNetwork := viper.GetString("network")
    cProjectId := viper.GetString("project-id")
    cProjectName := viper.GetString("project-name")
    cTelescopeUsername := viper.GetString("telescope-username")
    cTelescopePassword := viper.GetString("telescope-password")
    cRemoteWriteUrl := viper.GetString("remote-write-url")
    // cLogsSinkURL := viper.GetString("logs-sink-url")
    // cLokiUsername := viper.GetString("telescope-loki-username")
    // cLokiPassword := viper.GetString("telescope-loki-password")

    // Convert networksConfig.ScrapeConfig to local ScrapeConfig
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

    return Config{
        Server: ServerConfig{
            LogLevel: "info",
        },
        Metrics: MetricsConfig{
            Wal_Directory: "/tmp/telescope", // TODO: make this configurable
            Global: GlobalConfig{
                ScrapeInterval: "15s",
                ExternalLabels: map[string]string{
                    "project_id":   cProjectId,
                    "project_name": cProjectName,
                },
                RemoteWrite: []RemoteWrite{
                    {
                        URL: cRemoteWriteUrl,
                        BasicAuth: map[string]string{
                            "username": cTelescopeUsername,
                            "password": cTelescopePassword,
                        },
                    },
                },
            },
            Configs: []MetricConfig{
                {
                    Name:          toLowerAndEscape(cProjectName + "_" + cNetwork + "_metrics"),
                    HostFilter:    false,
                    ScrapeConfigs: scrapeConfigs,
                },
            },
        },
        Integrations: map[string]interface{}{
            "agent":         ToIntegrate{Enabled: false},
            "node_exporter": ToIntegrate{Enabled: true},
        },
    }
}

func writeConfigToFile(config Config, filePath string) error {
    var buf strings.Builder
    encoder := yaml.NewEncoder(&buf)
    encoder.SetIndent(2)

    if err := encoder.Encode(&config); err != nil {
        return fmt.Errorf("error marshaling to YAML: %v", err)
    }

    data := buf.String()
    // Optional: Perform any necessary string replacements or adjustments
    fixedData := strings.ReplaceAll(data, "job_name", "job_name")
    fixedData = strings.ReplaceAll(fixedData, "static_configs", "static_configs")

    return os.WriteFile(filePath, []byte(fixedData), 0644)
}

func networkDiscovery(network string) ([]string, error) {
    scrapeConfig, ok := networkConfigs[network]
    if !ok {
        return nil, fmt.Errorf("invalid network. Please choose from: ethereum, polkadot, hyperbridge, ssv")
    }
    return scrapeConfig.NetworkDiscovery()
}

func generateNetworkConfig() Config {
    cNetwork := viper.GetString("network")
    cProjectId := viper.GetString("project-id")
    cProjectName := viper.GetString("project-name")
    cTelescopeUsername := viper.GetString("telescope-username")
    cTelescopePassword := viper.GetString("telescope-password")
    cRemoteWriteUrl := viper.GetString("remote-write-url")
    cLogsSinkURL := viper.GetString("logs-sink-url")
    cLokiUsername := viper.GetString("telescope-loki-username")
    cLokiPassword := viper.GetString("telescope-loki-password")
    isEnableLogs := viper.GetBool("enable-logs")

    ports, err := networkDiscovery(cNetwork)
    handleErr(err, "Failed to discover blockchain port")

    for _, port := range ports {
        viper.Set("scrape_port", port)
    }

    networkConfig, ok := networkConfigs[cNetwork]
    if !ok {
        log.Fatalf("Invalid network configuration for: %v", cNetwork)
    }

    networkScrapeConfigs := networkConfig.GenerateScrapeConfigs(cProjectName, cNetwork)

    // Convert networksConfig.ScrapeConfig to local ScrapeConfig
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

    config := Config{
        Server: ServerConfig{
            LogLevel: "info",
        },
        Metrics: MetricsConfig{
            Wal_Directory: "/tmp/telescope", // TODO: make this configurable
            Global: GlobalConfig{
                ScrapeInterval: "15s",
                ExternalLabels: map[string]string{
                    "project_id":   cProjectId,
                    "project_name": cProjectName,
                },
                RemoteWrite: []RemoteWrite{
                    {
                        URL: cRemoteWriteUrl,
                        BasicAuth: map[string]string{
                            "username": cTelescopeUsername,
                            "password": cTelescopePassword,
                        },
                    },
                },
            },
            Configs: []MetricConfig{
                {
                    Name:          toLowerAndEscape(cProjectName + "_" + cNetwork + "_metrics"),
                    HostFilter:    false,
                    ScrapeConfigs: scrapeConfigs,
                },
            },
        },
        Integrations: map[string]interface{}{
            "agent":         ToIntegrate{Enabled: false},
            "node_exporter": ToIntegrate{Enabled: true},
        },
    }

    if isEnableLogs {
        // Include only Loki-specific LogClient
        config.Logs = LogsConfig{
            Configs: []LogConfig{
                {
                    Name: "telescope_logs",
                    Clients: []LogClient{
                        {
                            URL: cLogsSinkURL,
                            BasicAuth: BasicAuth{
                                Username: cLokiUsername,
                                Password: cLokiPassword,
                            },
                            ExternalLabels: map[string]string{
                                "project_id":   cProjectId,
                                "project_name": cProjectName,
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

    return config
}

func LoadNetworkConfig() string {
    config := generateNetworkConfig()

    // Create a custom marshaler
    var buf bytes.Buffer
    encoder := yaml.NewEncoder(&buf)
    encoder.SetIndent(2)

    if err := encoder.Encode(&config); err != nil {
        log.Fatalf("error marshaling to YAML: %v", err)
    }

    data := buf.String()
    // fixedData := strings.ReplaceAll(string(data), "job_name", "job_name")
    // fixedData = strings.ReplaceAll(fixedData, "static_configs", "static_configs")

    configFilePath := "telescope_config.yaml"
    if err := os.WriteFile(configFilePath, []byte(data), 0644); err != nil {
        log.Fatalf("error writing to file: %v", err)
    }

    return configFilePath
}

func checkRequiredFlags() error {
    requiredFlags := []string{
        "metrics",
        "enable-logs",
        "network",
        "project-id",
        "project-name",
        "telescope-username",
        "telescope-password",
        "remote-write-url",
        "telescope-loki-username", // New required flag
        "telescope-loki-password", // New required flag
    }
    missingFlags := []string{}

    for _, flag := range requiredFlags {
        if viper.GetString(flag) == "" && !viper.GetBool(flag) { // For bool flags
            missingFlags = append(missingFlags, flag)
        }
    }

    if len(missingFlags) > 0 && viper.GetString("config-file") == "" {
        return fmt.Errorf("error: missing required flags: %s", strings.Join(missingFlags, ", "))
    }

    return nil
}

func init() {
    prometheus.MustRegister(build.NewCollector("agent"))
    cobra.OnInitialize(initConfig)
    cmd.SetHelpFunc(helperFunction)
    cmd.Flags().StringVar(&cfgFile, "config-file", "", "Specify the config file")
    cmd.Flags().Bool("metrics", true, "Enable metrics")
    cmd.Flags().Bool("enable-logs", false, "Enable logs")
    cmd.Flags().String("network", "", "Specify the network")
    cmd.Flags().String("project-id", "", "Specify the project ID")
    cmd.Flags().String("project-name", "", "Specify the project name")
    cmd.Flags().String("telescope-username", "", "Specify the telescope username")
    cmd.Flags().String("telescope-password", "", "Specify the telescope password")
    cmd.Flags().String("remote-write-url", "", "Specify the remote write URL")
    cmd.Flags().String("logs-sink-url", "", "Specify the Log Sink URL")

    // Define new flags for Loki credentials
    cmd.Flags().String("telescope-loki-username", "", "Specify the Loki username")   // New flag
    cmd.Flags().String("telescope-loki-password", "", "Specify the Loki password")   // New flag

    // Bind flags with Viper
    viper.BindPFlag("config-file", cmd.Flags().Lookup("config-file"))
    viper.BindPFlag("metrics", cmd.Flags().Lookup("metrics"))
    viper.BindPFlag("enable-logs", cmd.Flags().Lookup("enable-logs"))
    viper.BindPFlag("network", cmd.Flags().Lookup("network"))
    viper.BindPFlag("project-id", cmd.Flags().Lookup("project-id"))
    viper.BindPFlag("project-name", cmd.Flags().Lookup("project-name"))
    viper.BindPFlag("telescope-username", cmd.Flags().Lookup("telescope-username"))
    viper.BindPFlag("telescope-password", cmd.Flags().Lookup("telescope-password"))
    viper.BindPFlag("remote-write-url", cmd.Flags().Lookup("remote-write-url"))
    viper.BindPFlag("logs-sink-url", cmd.Flags().Lookup("logs-sink-url"))

    // Bind new Loki flags
    viper.BindPFlag("telescope-loki-username", cmd.Flags().Lookup("telescope-loki-username")) // New binding
    viper.BindPFlag("telescope-loki-password", cmd.Flags().Lookup("telescope-loki-password")) // New binding
}

func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
        handleErr(viper.ReadInConfig(), "Failed to read config file")
    } else {
        viper.AutomaticEnv() // read in environment variables that match
    }
}

func (c *TelescopeConfig) loadConfig() error {
    if viper.GetString("config-file") != "" {
        cfgFile = viper.GetString("config-file")
        viper.SetConfigFile(cfgFile)
        if err := viper.ReadInConfig(); err != nil {
            return fmt.Errorf("Error reading config file: %s", err)
        }
        log.Println("Configuration loaded from file:", cfgFile)
    } else {
        LoadNetworkConfig()
        log.Println("Configuration generated and written to telescope_config.yaml")
    }

    // Check for required fields only if config-file is not provided
    if err := checkRequiredFlags(); err != nil {
        return err
    }

    // Assign the fields from Viper to the struct
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

    return nil
}

func agent(configPath string) {

    defaultCfg := server.DefaultConfig()
    logger := server.NewLogger(&defaultCfg)

    reloader := func(log *server.Logger) (*config.Config, error) {
        fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
        // fs.String("telescope", configPath, "Path to configuration file")
        // fs.Parse([]string{"--config.file=" + configPath}) // Ensure the flag is set
        return config.Load(fs, []string{"-config.file", configPath}, log)
    }
    cfg, err := reloader(logger)
    if err != nil {
        log.Fatalln(err)
    }

    // After this point we can start using go-kit logging.
    logger = server.NewLogger(cfg.Server)
    util_log.Logger = logger

    level.Info(logger).Log("boringcrypto enabled", boringcrypto.Enabled)
    ep, err := NewEntrypoint(logger, cfg, reloader)
    if err != nil {
        level.Error(logger).Log("msg", "error creating the agent server entrypoint", "err", err)
        os.Exit(1)
    }

    if err = ep.Start(); err != nil {
        level.Error(logger).Log("msg", "error running agent", "err", err)
        // Don't os.Exit here; we want to do cleanup by stopping promMetrics
    }

    ep.Stop()
    level.Info(logger).Log("msg", "agent exiting")
}

func main() {
	var config TelescopeConfig
	var scrapeConfigs []networksConfig.ScrapeConfig

	cobra.OnInitialize(initConfig)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	network = viper.GetString("network")
	projectName := viper.GetString("project-name")
	// if network == "" {
	// 	log.Fatal("Unsupported network: network argument is missing")
	// }
	networkConfig := getNetworkConfig(network)

	fmt.Println("Starting Telescope Agent", networkConfig)
	scrapeConfigs = networkConfig.GenerateScrapeConfigs(projectName, network)
	fullConfig := generateFullConfig(config, scrapeConfigs)
	configFilePath := "telescope_config.yaml"
	if err := writeConfigToFile(fullConfig, configFilePath); err != nil {
		log.Fatalf("Failed to write config to file: %v", err)
	}

	fmt.Printf("Configuration written to %s\n", configFilePath)

	agent(configFilePath)
}
