package build

import (
	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter/postgres"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/postgres_exporter"
	"github.com/grafana/river/rivertypes"
)

func (b *ConfigBuilder) appendPostgresExporter(config *postgres_exporter.Config, instanceKey *string) discovery.Exports {
	args := toPostgresExporter(config)
	return b.appendExporterBlock(args, config.Name(), instanceKey, "postgres")
}

func toPostgresExporter(config *postgres_exporter.Config) *postgres.Arguments {
	dataSourceNames := make([]rivertypes.Secret, 0)
	for _, dsn := range config.DataSourceNames {
		dataSourceNames = append(dataSourceNames, rivertypes.Secret(dsn))
	}

	return &postgres.Arguments{
		DataSourceNames:         dataSourceNames,
		DisableSettingsMetrics:  config.DisableSettingsMetrics,
		DisableDefaultMetrics:   config.DisableDefaultMetrics,
		CustomQueriesConfigPath: config.QueryPath,
		AutoDiscovery: postgres.AutoDiscovery{
			Enabled:           config.AutodiscoverDatabases,
			DatabaseAllowlist: config.IncludeDatabases,
			DatabaseDenylist:  config.ExcludeDatabases,
		},
	}
}
