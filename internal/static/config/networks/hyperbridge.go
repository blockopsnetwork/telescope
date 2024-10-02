package networks

import (
	"fmt"
)

type HyperbridgeConfig struct {
	Chain    string
	NodeType map[string]int
	Port     int
}

func NewHyperbridgeConfig() *HyperbridgeConfig {
	return &HyperbridgeConfig{
		Chain:    "hyperbridge",
		NodeType: map[string]int{"node": 8080},
		Port:     8080,
	}
}

func (h *HyperbridgeConfig) NetworkDiscovery() ([]string, error) {
	ports := []string{}
	for _, port := range h.NodeType {
		ports = append(ports, fmt.Sprintf("localhost:%d", port))
	}
	return ports, nil
}

func (h *HyperbridgeConfig) GenerateScrapeConfigs(projectName, network string) []ScrapeConfig {
	var scrapeConfigs []ScrapeConfig
	idx := 0
	for nodeType, port := range h.NodeType {
		jobName := fmt.Sprintf("%s_%s_%s_job_%d", toLowerAndEscape(projectName), network, nodeType, idx)
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
	return scrapeConfigs
}

func (h *HyperbridgeConfig) AutoconfigureScrapeConfigs(projectName, network string) ([]ScrapeConfig, error) {
	// Implement the logic for auto-configuring scrape configs
	return h.GenerateScrapeConfigs(projectName, network), nil
}