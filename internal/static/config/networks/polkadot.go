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
		NodeType: map[string]int{"relaychain": 9615, "parachain": 9616},
		Port:     9615,
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