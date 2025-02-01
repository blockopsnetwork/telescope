package build

import (
	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter/oracledb"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/oracledb_exporter"
	"github.com/grafana/river/rivertypes"
)

func (b *ConfigBuilder) appendOracledbExporter(config *oracledb_exporter.Config, instanceKey *string) discovery.Exports {
	args := toOracledbExporter(config)
	return b.appendExporterBlock(args, config.Name(), instanceKey, "oracledb")
}

func toOracledbExporter(config *oracledb_exporter.Config) *oracledb.Arguments {
	return &oracledb.Arguments{
		ConnectionString: rivertypes.Secret(config.ConnectionString),
		MaxIdleConns:     config.MaxIdleConns,
		MaxOpenConns:     config.MaxOpenConns,
		QueryTimeout:     config.QueryTimeout,
	}
}
