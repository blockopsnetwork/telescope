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
		Chain:    "gargantua",
		NodeType: map[string]int{"relaychain": 9615, "parachain": 9616},
		Port:     9615,
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