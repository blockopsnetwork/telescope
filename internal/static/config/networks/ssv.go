package networks

import (
	"fmt"
	"github.com/ethpandaops/beacon/pkg/beacon"
	"github.com/ethpandaops/ethereum-metrics-exporter/pkg/exporter/execution"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

)

type NodeConfig struct {
	Type string
	Port int
}

type SSVConfig struct {
	Protocol  string
	NodeTypes []NodeConfig
}

func NewSSVConfig() *SSVConfig {
	return &SSVConfig{
		Protocol: "ssv",
		NodeTypes: []NodeConfig{
			{Type: "execution", Port: 6060},
			{Type: "consensus", Port: 8008},
			{Type: "mevboost", Port: 18550},
			{Type: "ssvdkg", Port: 3030},
			{Type: "ssv", Port: 13000},
		},
	}
}

// Update the method signature to match the expected interface
func (s *SSVConfig) GenerateScrapeConfigs(projectName, protocol string) []ScrapeConfig {
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

func (s *SSVConfig) NetworkDiscovery() ([]string, error) {
	var ports []string
	
	for _, node := range s.NodeTypes { 
		ports = append(ports, fmt.Sprintf("localhost:%d", node.Port))
	}
	
	return ports, nil
}


