package main

import (
	"flag"
	"log"
	"os"
	"fmt"
	"strings"
	"net/url"


	"github.com/go-kit/log/level"
	"gopkg.in/yaml.v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/grafana/agent/internal/boringcrypto"
	"github.com/grafana/agent/internal/build"
	"github.com/grafana/agent/internal/static/config"
	"github.com/grafana/agent/internal/static/server"
	util_log "github.com/grafana/agent/internal/util/log"

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

var cfgFile string

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Metrics MetricsConfig `yaml:"metrics"`
	Logs        LogsConfig        `yaml:"logs"`
	Integrations       map[string]interface{}  `yaml:"integrations"`
	
}

type ServerConfig struct {
	LogLevel string `yaml:"log_level"`
}

type MetricsConfig struct {
	Global   GlobalConfig   `yaml:"global"`
	Wal_Directory string `yaml:"wal_directory"`
	Configs  []MetricConfig `yaml:"configs"`
}

type GlobalConfig struct {
	ScrapeInterval string `yaml:"scrape_interval"`
	ExternalLabels      map[string]string  `yaml:"external_labels"`
	RemoteWrite    []RemoteWrite `yaml:"remote_write"`
}

type MetricConfig struct {
	Name           string         `yaml:"name"`
	HostFilter     bool           `yaml:"host_filter"`
	ScrapeConfigs  []ScrapeConfig `yaml:"scrape_configs"`
}

type ScrapeConfig struct {
	JobName       string          `yaml:"job_name"`
	StaticConfigs []StaticConfig  `yaml:"static_configs"`
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
	Targets []string           `yaml:"targets"`
}

type RemoteWrite struct {
	URL string `yaml:"url"`
	BasicAuth      map[string]string  `yaml:"basic_auth"`

}

type LogsConfig struct {
	Configs []LogConfig `yaml:"configs"`
}

type LogConfig struct {
	Name      string     `yaml:"name"`
	Clients   []LogClient `yaml:"clients"`
	Positions Positions  `yaml:"positions"`
}

type Positions struct {
	Filename string `yaml:"filename"`
}

type LogClient struct {
	URL           string            `yaml:"url"`
	BasicAuth     BasicAuth         `yaml:"basic_auth"`
	ExternalLabels map[string]string `yaml:"external_labels"`
}


type TelescopeConfig struct {
	Metrics bool
	Logs bool
	Network string
	ProjectId string
	ProjectName string
	TelescopeUsername string
	TelescopePassword string
	RemoteWriteUrl string
}

type BlockchainNetworkConfig struct {
	Chain string //or env e.g testnet
	Network string
	NodeType map[string]int
	Port int
}

var networkConfigs = map[string]BlockchainNetworkConfig{
	"ethereum": {
		Chain:    "sepolia",
		Network:  "ethereum",
		NodeType: map[string]int{"execution": 6060, "consensus": 8008},
		Port:     6060,
	},
	"polkadot": {
		Chain:    "polkadot",
		Network:  "polkadot",
		NodeType: map[string]int{"relaychain": 9615, "parachain": 9616},
		Port:     9615,
	},
	"hyperbridge": {
		Chain:    "gargantua",
		Network:  "hyperbridge",
		NodeType: map[string]int{"relaychain": 9615, "parachain": 9616},
		Port:     9615,
	},
}

func toLowerAndEscape(input string) string {
	lowercase := strings.ToLower(input) // Convert to lowercase
	escaped := url.QueryEscape(lowercase) // Escape special characters for URL
	return escaped
}


func networkDiscovery(network string) ([]string, error) {
	scrapeConfig, ok := networkConfigs[network]

	if !ok {
		return nil, fmt.Errorf("invalid network. Please choose from: ethereum, polkadot, hyperbridge")
	}

	ports := []string{}
	for _, port := range scrapeConfig.NodeType {
		ports = append(ports, fmt.Sprintf("localhost:%d", port))
	}

	if len(ports) == 0 {
		return []string{fmt.Sprintf("default port: %d", scrapeConfig.Port)}, nil
	}

	return ports, nil
}


var cmd = &cobra.Command{
	Use:   "metrics",
	Short: "A Web3 Observability tooling",
	Long: `Gain insights into your Web3 infrastructure with Telescope.`,
	Run: func(cmd *cobra.Command, args []string) {
		var config TelescopeConfig
		if err := config.loadConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
	
}


func generateNetworkConfig() Config {
	// var telescope_config TelescopeConfig
	// cMetrics := viper.GetBool("metrics")
	cNetwork := viper.GetString("network")
	cProjectId := viper.GetString("project-id")
	cProjectName := viper.GetString("project-name")
	cTelescopeUsername := viper.GetString("telescope-username")
	cTelescopePassword := viper.GetString("telescope-password")
	cRemoteWriteUrl := viper.GetString("remote-write-url")
	isEnableLogs := viper.GetString("enable-logs")
	cLogSinkURL := viper.GetString("logs-sink-url")

	ports, err := networkDiscovery(cNetwork)
	

	if err != nil {
		log.Fatalf("Unable to discover blockchain port: %v", err)
	}

	for _, port := range ports {
		viper.Set("scrape_port", port)
	}

	networkConfig, ok := networkConfigs[cNetwork]
    if !ok {
        log.Fatalf("Invalid network configuration for: %v", cNetwork)
    }

	var scrapeConfigs []ScrapeConfig
    idx := 0  // Initialize index for job naming
    for nodeType, port := range networkConfig.NodeType {
        jobName := fmt.Sprintf("%s_%s_%s_job_%d", toLowerAndEscape(cProjectName), cNetwork, nodeType, idx)
        target := fmt.Sprintf("localhost:%d", port)
        scrapeConfigs = append(scrapeConfigs, ScrapeConfig{
            JobName: jobName,
            StaticConfigs: []StaticConfig{
                {
                    Targets: []string{target},
                },
            },
        })
        idx++
    }

	config := Config{
		Server: ServerConfig{
			LogLevel: "info",
		},
		Metrics: MetricsConfig{
			Wal_Directory: "/tmp/wal",
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
					Name:        toLowerAndEscape(cProjectName + cNetwork + "_metrics"),
					HostFilter:  false,
					ScrapeConfigs: scrapeConfigs,
				},
			},
		},
		Integrations: map[string]interface{}{
			"agent": ToIntegrate{
				Enabled: false,
			},
			"node_exporter": ToIntegrate{
				Enabled: true,
			},
		},
	}

	if isEnableLogs == "true" {
		config.Logs = LogsConfig{
			Configs: []LogConfig{
				{
					Name: "telescope_logs",
					Clients: []LogClient{
						{
							URL: cLogSinkURL,
							BasicAuth: BasicAuth{
								Username: cTelescopeUsername,
								Password: cTelescopePassword,
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
	config := generateNetworkConfig() // Assuming this function returns a `Config` object correctly populated

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error marshaling to YAML: %v", err)
	}

	configFilePath := "telescope_config.yaml"
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		log.Fatalf("error writing to file: %v", err)
	}

	return configFilePath
}



var helperFunction = func(cmd *cobra.Command, args []string) {
	fmt.Println("Custom Help:")
    fmt.Println("Usage: app start [flags]")
    fmt.Println("\nFlags:")
    fmt.Println("  --metrics\t\tEnable metrics")
    fmt.Println("  --network\t\tSpecify the network")
    fmt.Println("  --project-id\t\tSpecify the project ID")
    fmt.Println("  --project-name\tSpecify the project name")
    fmt.Println("  --telescope-username\tSpecify the telescope username")
    fmt.Println("  --telescope-password\tSpecify the telescope password")
    fmt.Println("  --remote-write-url\tSpecify the remote write URL")
	fmt.Println("  --config-file\t\tSpecify the config file")
	fmt.Println("  --logs\t\tEnable logs")
}


func checkRequiredFlags() error {
    requiredFlags := []string{"metrics", "enable-logs", "network", "project-id", "project-name", "telescope-username", "telescope-password", "remote-write-url"}
    missingFlags := []string{}

    for _, flag := range requiredFlags {
        if viper.GetString(flag) == "" {
            missingFlags = append(missingFlags, flag)
        }
    }

    if len(missingFlags) > 0 && viper.GetString("config-file") == "" {
        return fmt.Errorf("Error: Missing required flags: %s", strings.Join(missingFlags, ", "))
    }

    return nil
}

func init() {
	prometheus.MustRegister(build.NewCollector("agent"))
	cobra.OnInitialize(initConfig)
	// cmd.SetHelpFunc(helperFunction)
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

    // Bind flags with viper
	viper.BindPFlag("config-file", cmd.Flags().Lookup("config-file"))
    viper.BindPFlag("metrics", cmd.Flags().Lookup("metrics"))
    viper.BindPFlag("network", cmd.Flags().Lookup("network"))
    viper.BindPFlag("project-id", cmd.Flags().Lookup("project-id"))
    viper.BindPFlag("project-name", cmd.Flags().Lookup("project-name"))
    viper.BindPFlag("telescope-username", cmd.Flags().Lookup("telescope-username"))
    viper.BindPFlag("telescope-password", cmd.Flags().Lookup("telescope-password"))
    viper.BindPFlag("remote-write-url", cmd.Flags().Lookup("remote-write-url"))
	viper.BindPFlag("logs-sink-url", cmd.Flags().Lookup("logs-sink-url"))
	viper.BindPFlag("enable-logs", cmd.Flags().Lookup("enable-logs"))

	
	// cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
	} else {
		viper.AutomaticEnv() // read in environment variables that match
	}
}

func (c *TelescopeConfig) loadConfig() error {

	// Check if config-file is provided
	if viper.GetString("config-file") != "" {
		cfgFile = viper.GetString("config-file")
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("Error reading config file: %s", err)
		}
	} else {
		// generateNetworkConfig 
		LoadNetworkConfig()
	}

	// Check for required fields only if config-file is not provided
	if err := checkRequiredFlags(); err != nil {
		return err
	}


	return nil


}


// func init() {
// 	prometheus.MustRegister(build.NewCollector("agent"))
// }

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
	
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	var configPath string
	if cfgFile != "" {
		configPath = cfgFile
	} else {
		configPath = LoadNetworkConfig()
	}

	agent(configPath)
}
