package build

import (
	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter/squid"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/squid_exporter"
	"github.com/grafana/river/rivertypes"
)

func (b *ConfigBuilder) appendSquidExporter(config *squid_exporter.Config, instanceKey *string) discovery.Exports {
	args := toSquidExporter(config)
	return b.appendExporterBlock(args, config.Name(), instanceKey, "squid")
}

func toSquidExporter(config *squid_exporter.Config) *squid.Arguments {
	return &squid.Arguments{
		SquidAddr:     config.Address,
		SquidUser:     config.Username,
		SquidPassword: rivertypes.Secret(config.Password),
	}
}
