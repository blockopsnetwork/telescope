package networks

import (
	"fmt"
)

type SSVConfig struct {
	Chain    string
	NodeType map[string]int
	Port     int
}

func NewSSVConfig() *SSVConfig {
	return &SSVConfig{
		Chain: "ssv",
		NodeType: map[string]int{
			"execution": 6060,
			"consensus": 8008,
			"ssv":       15000, // Assuming SSV uses port 15000, adjust as needed
		},
		Port: 15000,
	}
}

func (s *SSVConfig) NetworkDiscovery() ([]string, error) {
	ports := []string{}
	for _, port := range s.NodeType {
		ports = append(ports, fmt.Sprintf("localhost:%d", port))
	}
	return ports, nil
}

func (s *SSVConfig) GenerateScrapeConfigs(projectName, network string) []ScrapeConfig {
	var scrapeConfigs []ScrapeConfig
	idx := 0
	for nodeType, port := range s.NodeType {
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