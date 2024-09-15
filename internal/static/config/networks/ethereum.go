package networks

import (
	"fmt"
)

type EthereumConfig struct {
	Chain    string
	NodeType map[string]int
	Port     int
}

func NewEthereumConfig() *EthereumConfig {
	return &EthereumConfig{
		Chain:    "sepolia",
		NodeType: map[string]int{"execution": 6060, "consensus": 8008},
		Port:     6060,
	}
}

func (e *EthereumConfig) NetworkDiscovery() ([]string, error) {
	ports := []string{}
	for _, port := range e.NodeType {
		ports = append(ports, fmt.Sprintf("localhost:%d", port))
	}
	return ports, nil
}

func (e *EthereumConfig) GenerateScrapeConfigs(projectName, network string) []ScrapeConfig {
	var scrapeConfigs []ScrapeConfig
	idx := 0
	for nodeType, port := range e.NodeType {
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