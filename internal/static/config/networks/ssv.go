package networks

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type NodeConfig struct {
	Type string
	Port int
}

type SSVConfig struct {
	Chain           string
	ExecutionNodes  []NodeConfig
	ConsensusNodes  []NodeConfig
	ValidatorNodes  []NodeConfig
	SSVNode         NodeConfig
	MEVBoost        NodeConfig
	SSVDKG          NodeConfig
	DefaultSSVPort  int
}

type NodeDiscovery struct {
	Type       string
	Port       int
	APIEndpoint string
}

func NewSSVConfig() *SSVConfig {
	return &SSVConfig{
		Chain: "ssv",
		ExecutionNodes: []NodeConfig{
			{Type: "geth", Port: 6060},
			{Type: "nethermind", Port: 6060},
			{Type: "besu", Port: 9545},
			{Type: "erigon", Port: 6060},
			{Type: "reth", Port: 9090},
		},
		ConsensusNodes: []NodeConfig{
			{Type: "prysm", Port: 8080},
			{Type: "lighthouse", Port: 5054},
			{Type: "teku", Port: 8008},
			{Type: "nimbus", Port: 8008},
			{Type: "lodestar", Port: 8008},
		},
		ValidatorNodes: []NodeConfig{
			{Type: "prysm", Port: 8081},
			{Type: "lighthouse", Port: 5064},
			{Type: "teku", Port: 8008},  
			{Type: "nimbus", Port: 8008}, 
			{Type: "lodestar", Port: 5064},
		},
		SSVNode:        NodeConfig{Type: "ssv", Port: 13000},
		MEVBoost:       NodeConfig{Type: "mevboost", Port: 18550},
		SSVDKG:         NodeConfig{Type: "ssvdkg", Port: 3030},
		DefaultSSVPort: 13000,
	}
}

func (s *SSVConfig) AutoconfigureScrapeConfigs(projectName, network string) ([]ScrapeConfig, error) {
	discoveredNodes, err := AutodiscoverNodes()
	if err != nil {
		return nil, err
	}

	var scrapeConfigs []ScrapeConfig

	for _, node := range discoveredNodes {
		scrapeConfigs = append(scrapeConfigs, ScrapeConfig{
			JobName: fmt.Sprintf("%s_%s_%s", toLowerAndEscape(projectName), network, node.Type),
			StaticConfigs: []StaticConfig{
				{
					Targets: []string{fmt.Sprintf("localhost:%d", node.Port)},
				},
			},
		})
	}

	return scrapeConfigs, nil
}

func (s *SSVConfig) NetworkDiscovery() ([]string, error) {
	var ports []string
	
	for _, node := range s.ExecutionNodes {
		ports = append(ports, fmt.Sprintf("localhost:%d", node.Port))
	}
	for _, node := range s.ConsensusNodes {
		ports = append(ports, fmt.Sprintf("localhost:%d", node.Port))
	}
	for _, node := range s.ValidatorNodes {
		ports = append(ports, fmt.Sprintf("localhost:%d", node.Port))
	}
	ports = append(ports, 
		fmt.Sprintf("localhost:%d", s.SSVNode.Port),
		fmt.Sprintf("localhost:%d", s.MEVBoost.Port),
		fmt.Sprintf("localhost:%d", s.SSVDKG.Port),
	)
	
	return ports, nil
}

func (s *SSVConfig) GenerateScrapeConfigs(projectName, network string) []ScrapeConfig {
	var scrapeConfigs []ScrapeConfig
	
	// Execution nodes
	for _, node := range s.ExecutionNodes {
		scrapeConfigs = append(scrapeConfigs, ScrapeConfig{
			JobName: fmt.Sprintf("%s_%s_execution_%s", toLowerAndEscape(projectName), network, node.Type),
			StaticConfigs: []StaticConfig{
				{
					Targets: []string{fmt.Sprintf("localhost:%d", node.Port)},
				},
			},
		})
	}
	
	// Consensus nodes
	for _, node := range s.ConsensusNodes {
		scrapeConfigs = append(scrapeConfigs, ScrapeConfig{
			JobName: fmt.Sprintf("%s_%s_consensus_%s", toLowerAndEscape(projectName), network, node.Type),
			StaticConfigs: []StaticConfig{
				{
					Targets: []string{fmt.Sprintf("localhost:%d", node.Port)},
				},
			},
		})
	}
	
	// Validator nodes
	for _, node := range s.ValidatorNodes {
		scrapeConfigs = append(scrapeConfigs, ScrapeConfig{
			JobName: fmt.Sprintf("%s_%s_validator_%s", toLowerAndEscape(projectName), network, node.Type),
			StaticConfigs: []StaticConfig{
				{
					Targets: []string{fmt.Sprintf("localhost:%d", node.Port)},
				},
			},
		})
	}
	
	// SSV Node, MEV-Boost, and SSV DKG
	for _, node := range []NodeConfig{s.SSVNode, s.MEVBoost, s.SSVDKG} {
		scrapeConfigs = append(scrapeConfigs, ScrapeConfig{
			JobName: fmt.Sprintf("%s_%s_%s", toLowerAndEscape(projectName), network, node.Type),
			StaticConfigs: []StaticConfig{
				{
					Targets: []string{fmt.Sprintf("localhost:%d", node.Port)},
				},
			},
		})
	}

	return scrapeConfigs
}


func AutodiscoverNodes() ([]NodeConfig, error) {
	nodesToCheck := []NodeDiscovery{
		{"geth", 6060, "/metrics"},
		{"nethermind", 6060, "/metrics"},
		{"besu", 9545, "/metrics"},
		{"erigon", 6060, "/metrics"},
		{"reth", 9090, "/metrics"},
		{"prysm_beacon", 8080, "/metrics"},
		{"lighthouse_beacon", 5054, "/metrics"},
		{"teku", 8008, "/metrics"},
		{"nimbus", 8008, "/metrics"},
		{"lodestar", 8008, "/metrics"},
		{"prysm_validator", 8081, "/metrics"},
		{"lighthouse_validator", 5064, "/metrics"},
		{"lodestar_validator", 5064, "/metrics"},
		{"ssv", 13000, "/metrics"},
		{"mevboost", 18550, "/metrics"},
		{"ssvdkg", 3030, "/metrics"},
	}

	var discoveredNodes []NodeConfig

	for _, node := range nodesToCheck {
		if isPortOpen(node.Port) && isAPIAccessible(fmt.Sprintf("http://localhost:%d%s", node.Port, node.APIEndpoint)) {
			discoveredNodes = append(discoveredNodes, NodeConfig{Type: node.Type, Port: node.Port})
		}
	}

	return discoveredNodes, nil
}

func isPortOpen(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func isAPIAccessible(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

