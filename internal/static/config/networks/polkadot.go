package networks

import (
	"fmt"
)

type PolkadotConfig struct {
	Chain    string
	NodeType map[string]int
	Port     int
}

func NewPolkadotConfig() *PolkadotConfig {
	return &PolkadotConfig{
		Chain:    "polkadot",
		NodeType: map[string]int{"relaychain": 30333, "parachains": 9933},
		Port:     30333,
	}
}

func (p *PolkadotConfig) NetworkDiscovery() ([]string, error) {
	ports := []string{}
	for _, port := range p.NodeType {
		ports = append(ports, fmt.Sprintf("localhost:%d", port))
	}
	return ports, nil
}

func (p *PolkadotConfig) GenerateScrapeConfigs(projectName, network string) []ScrapeConfig {
	var scrapeConfigs []ScrapeConfig
	idx := 0
	for nodeType, port := range p.NodeType {
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

func (p *PolkadotConfig) AutoconfigureScrapeConfigs(projectName, network string) ([]ScrapeConfig, error) {
	// Implement the logic for auto-configuring scrape configs
	return p.GenerateScrapeConfigs(projectName, network), nil
}