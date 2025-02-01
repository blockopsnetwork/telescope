package build

import (
	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter/github"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/github_exporter"
	"github.com/grafana/river/rivertypes"
)

func (b *ConfigBuilder) appendGithubExporter(config *github_exporter.Config, instanceKey *string) discovery.Exports {
	args := toGithubExporter(config)
	return b.appendExporterBlock(args, config.Name(), instanceKey, "github")
}

func toGithubExporter(config *github_exporter.Config) *github.Arguments {
	return &github.Arguments{
		APIURL:        config.APIURL,
		Repositories:  config.Repositories,
		Organizations: config.Organizations,
		Users:         config.Users,
		APIToken:      rivertypes.Secret(config.APIToken),
		APITokenFile:  config.APITokenFile,
	}
}
