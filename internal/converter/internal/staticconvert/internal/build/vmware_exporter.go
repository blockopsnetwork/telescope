package build

import (
	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter/vsphere"
	vmware_exporter_v2 "github.com/blockopsnetwork/telescope/internal/static/integrations/v2/vmware_exporter"
	"github.com/grafana/river/rivertypes"
)

func (b *ConfigBuilder) appendVmwareExporterV2(config *vmware_exporter_v2.Config) discovery.Exports {
	args := toVmwareExporter(config)
	return b.appendExporterBlock(args, config.Name(), nil, "vsphere")
}

func toVmwareExporter(config *vmware_exporter_v2.Config) *vsphere.Arguments {
	return &vsphere.Arguments{
		ChunkSize:               config.ChunkSize,
		CollectConcurrency:      config.CollectConcurrency,
		VSphereURL:              config.VSphereURL,
		VSphereUser:             config.VSphereUser,
		VSpherePass:             rivertypes.Secret(config.VSpherePass),
		ObjectDiscoveryInterval: config.ObjectDiscoveryInterval,
		EnableExporterMetrics:   config.EnableExporterMetrics,
	}
}
