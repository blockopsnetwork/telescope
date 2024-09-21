package networks

type NetworkConfig interface {
    NetworkDiscovery() ([]string, error)
    GenerateScrapeConfigs(projectName, network string) []ScrapeConfig
    AutoconfigureScrapeConfigs(projectName, network string) ([]ScrapeConfig, error)
}