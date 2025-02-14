// Package zipkin provides an otelcol.receiver.zipkin component.
package zipkin

import (
	"github.com/blockopsnetwork/telescope/internal/component"
	"github.com/blockopsnetwork/telescope/internal/component/otelcol"
	"github.com/blockopsnetwork/telescope/internal/component/otelcol/receiver"
	"github.com/blockopsnetwork/telescope/internal/featuregate"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver"
	otelcomponent "go.opentelemetry.io/collector/component"
	otelextension "go.opentelemetry.io/collector/extension"
)

func init() {
	component.Register(component.Registration{
		Name:      "otelcol.receiver.zipkin",
		Stability: featuregate.StabilityStable,
		Args:      Arguments{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			fact := zipkinreceiver.NewFactory()
			return receiver.New(opts, fact, args.(Arguments))
		},
	})
}

// Arguments configures the otelcol.receiver.zipkin component.
type Arguments struct {
	ParseStringTags bool `river:"parse_string_tags,attr,optional"`

	HTTPServer otelcol.HTTPServerArguments `river:",squash"`

	// DebugMetrics configures component internal metrics. Optional.
	DebugMetrics otelcol.DebugMetricsArguments `river:"debug_metrics,block,optional"`

	// Output configures where to send received data. Required.
	Output *otelcol.ConsumerArguments `river:"output,block"`
}

var _ receiver.Arguments = Arguments{}

// SetToDefault implements river.Defaulter.
func (args *Arguments) SetToDefault() {
	*args = Arguments{
		HTTPServer: otelcol.HTTPServerArguments{
			Endpoint: "0.0.0.0:9411",
		},
	}
	args.DebugMetrics.SetToDefault()
}

// Convert implements receiver.Arguments.
func (args Arguments) Convert() (otelcomponent.Config, error) {
	return &zipkinreceiver.Config{
		ParseStringTags: args.ParseStringTags,
		ServerConfig:    *args.HTTPServer.Convert(),
	}, nil
}

// Extensions implements receiver.Arguments.
func (args Arguments) Extensions() map[otelcomponent.ID]otelextension.Extension {
	return nil
}

// Exporters implements receiver.Arguments.
func (args Arguments) Exporters() map[otelcomponent.DataType]map[otelcomponent.ID]otelcomponent.Component {
	return nil
}

// NextConsumers implements receiver.Arguments.
func (args Arguments) NextConsumers() *otelcol.ConsumerArguments {
	return args.Output
}

// DebugMetricsConfig implements receiver.Arguments.
func (args Arguments) DebugMetricsConfig() otelcol.DebugMetricsArguments {
	return args.DebugMetrics
}
