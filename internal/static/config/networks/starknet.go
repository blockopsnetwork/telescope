package networks

import (
	"fmt"
)

type StarknetNodeConfig struct {
	Type string
	Port int
}

type StarknetConfig struct {
	Protocol  string
	NodeTypes []StarknetNodeConfig
}

func NewStarknetConfig() *StarknetConfig {
	return &StarknetConfig{
		Protocol: "starknet",
		NodeTypes: []StarknetNodeConfig{
			{Type: "execution", Port: 6060},
			{Type: "consensus", Port: 8008},
			{Type: "juno", Port: 9090},
			{Type: "pathfinder", Port: 9090},
		},
	}
}

// Update the method signature to match the expected interface
func (s *StarknetConfig) GenerateScrapeConfigs(projectName, protocol string) []ScrapeConfig {
	var scrapeConfigs []ScrapeConfig

	for _, node := range s.NodeTypes {
		scrapeConfigs = append(scrapeConfigs, ScrapeConfig{
			JobName: fmt.Sprintf("%s_client_%s", protocol, node.Type ),
			StaticConfigs: []StaticConfig{
				{
					Targets: []string{fmt.Sprintf("localhost:%d", node.Port)},
				},
			},
		})
	}

	return scrapeConfigs
}

func (s *StarknetConfig) NetworkDiscovery() ([]string, error) {
	var ports []string

	for _, node := range s.NodeTypes {
		ports = append(ports, fmt.Sprintf("localhost:%d", node.Port))
	}

	return ports, nil
}


