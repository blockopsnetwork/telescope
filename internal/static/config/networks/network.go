package networks

type NetworkConfig interface {
    NetworkDiscovery() ([]string, error)
    GenerateScrapeConfigs(projectName, network string) []ScrapeConfig
}