package telescope

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/drone/envsubst/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gopkg.in/yaml.v3"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Metrics  MetricsConfig  `yaml:"metrics"`
    Configs  []interface{}  `yaml:"configs"` 
    Integrations IntegrationsConfig `yaml:"integrations"`
}

type ServerConfig struct {
    LogLevel string `yaml:"log_level"`
}

type MetricsConfig struct {
    WalDirectory string           `yaml:"wal_directory"`
    Global       GlobalMetricsConfig `yaml:"global"`
}

type GlobalMetricsConfig struct {
    ScrapeInterval string            `yaml:"scrape_interval"`
    ExternalLabels map[string]string `yaml:"external_labels"`
    RemoteWrite    []RemoteWriteConfig `yaml:"remote_write"`
}

type RemoteWriteConfig struct {
    URL       string `yaml:"url"`
    BasicAuth BasicAuthConfig `yaml:"basic_auth"`
}

type BasicAuthConfig struct {
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}

type IntegrationsConfig struct {
    Agent        AgentConfig        `yaml:"agent"`
    NodeExporter NodeExporterConfig `yaml:"node_exporter"`
}

type AgentConfig struct {
    Enabled bool `yaml:"enabled"`
}

type NodeExporterConfig struct {
    Enabled bool `yaml:"enabled"`
}

func PrepareConfig(config *Config) Config {
	return Config{
		Server: ServerConfig{
			LogLevel: "info",
		},
		Metrics: MetricsConfig{
			WalDirectory: "/tmp/wal",
			Global: GlobalMetricsConfig{
				ScrapeInterval: "15s",
				ExternalLabels: map[string]string{
					"project_id":   config.ProjectId,
                    "project_name": config.ProjectName,
				},
				RemoteWrite: []RemoteWriteConfig{
					{
						URL: config.RemoteWriteUrl,
						BasicAuth: BasicAuthConfig{
							Username: config.TelescopeUsername,
							Password: config.TelescopePassword,
						},
					},
				},
			},
		},
		
		// Blcockhain Network specofoc config
		Integrations: IntegrationsConfig{
			Agent: AgentConfig{
				Enabled: true,
			},
			NodeExporter: NodeExporterConfig{
				Enabled: true,
			},
		},
	}

}

func GenerateConfigFile(config *Config) error {
	cfg := PrepareConfig(config)

	// Marshal the config to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	// Write the config to a file
	telescopeConfigFile := "telescope/telescope.yaml"
	if err := os.WriteFile(telescopeConfigFile, data, 0644); err != nil {
        return fmt.Errorf("Error writing configuration file: %v", err)
    }

	fmt.Println("Configuration file created successfully.")
    return nil
}