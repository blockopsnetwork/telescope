package build

import (
	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter/dnsmasq"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/dnsmasq_exporter"
)

func (b *ConfigBuilder) appendDnsmasqExporter(config *dnsmasq_exporter.Config, instanceKey *string) discovery.Exports {
	args := toDnsmasqExporter(config)
	return b.appendExporterBlock(args, config.Name(), instanceKey, "dnsmasq")
}

func toDnsmasqExporter(config *dnsmasq_exporter.Config) *dnsmasq.Arguments {
	return &dnsmasq.Arguments{
		Address:      config.DnsmasqAddress,
		LeasesFile:   config.LeasesPath,
		ExposeLeases: config.ExposeLeases,
	}
}
