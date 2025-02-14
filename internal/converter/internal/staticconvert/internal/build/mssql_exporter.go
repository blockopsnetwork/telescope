package build

import (
	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter/mssql"
	mssql_exporter "github.com/blockopsnetwork/telescope/internal/static/integrations/mssql"
	"github.com/grafana/river/rivertypes"
)

func (b *ConfigBuilder) appendMssqlExporter(config *mssql_exporter.Config, instanceKey *string) discovery.Exports {
	args := toMssqlExporter(config)
	return b.appendExporterBlock(args, config.Name(), instanceKey, "mssql")
}

func toMssqlExporter(config *mssql_exporter.Config) *mssql.Arguments {
	return &mssql.Arguments{
		ConnectionString:   rivertypes.Secret(config.ConnectionString),
		MaxIdleConnections: config.MaxIdleConnections,
		MaxOpenConnections: config.MaxOpenConnections,
		Timeout:            config.Timeout,
	}
}
