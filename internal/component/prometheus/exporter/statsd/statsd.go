package statsd

import (
	"github.com/blockopsnetwork/telescope/internal/component"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter"
	"github.com/blockopsnetwork/telescope/internal/featuregate"
	"github.com/blockopsnetwork/telescope/internal/static/integrations"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.exporter.statsd",
		Stability: featuregate.StabilityStable,
		Args:      Arguments{},
		Exports:   exporter.Exports{},

		Build: exporter.New(createExporter, "statsd"),
	})
}

func createExporter(opts component.Options, args component.Arguments, defaultInstanceKey string) (integrations.Integration, string, error) {
	a := args.(Arguments)
	cfg, err := a.Convert()
	if err != nil {
		return nil, "", err
	}
	return integrations.NewIntegrationWithInstanceKey(opts.Logger, cfg, defaultInstanceKey)
}
