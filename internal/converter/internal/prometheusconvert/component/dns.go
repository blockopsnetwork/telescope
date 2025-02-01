package component

import (
	"time"

	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/discovery/dns"
	"github.com/blockopsnetwork/telescope/internal/converter/diag"
	"github.com/blockopsnetwork/telescope/internal/converter/internal/common"
	"github.com/blockopsnetwork/telescope/internal/converter/internal/prometheusconvert/build"
	prom_dns "github.com/prometheus/prometheus/discovery/dns"
)

func appendDiscoveryDns(pb *build.PrometheusBlocks, label string, sdConfig *prom_dns.SDConfig) discovery.Exports {
	discoveryDnsArgs := toDiscoveryDns(sdConfig)
	name := []string{"discovery", "dns"}
	block := common.NewBlockWithOverride(name, label, discoveryDnsArgs)
	pb.DiscoveryBlocks = append(pb.DiscoveryBlocks, build.NewPrometheusBlock(block, name, label, "", ""))
	return common.NewDiscoveryExports("discovery.dns." + label + ".targets")
}

func ValidateDiscoveryDns(sdConfig *prom_dns.SDConfig) diag.Diagnostics {
	return make(diag.Diagnostics, 0)
}

func toDiscoveryDns(sdConfig *prom_dns.SDConfig) *dns.Arguments {
	if sdConfig == nil {
		return nil
	}

	return &dns.Arguments{
		Names:           sdConfig.Names,
		RefreshInterval: time.Duration(sdConfig.RefreshInterval),
		Type:            sdConfig.Type,
		Port:            sdConfig.Port,
	}
}
