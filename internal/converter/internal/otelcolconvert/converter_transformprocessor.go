package otelcolconvert

import (
	"fmt"

	"github.com/blockopsnetwork/telescope/internal/component/otelcol"
	"github.com/blockopsnetwork/telescope/internal/component/otelcol/processor/transform"
	"github.com/blockopsnetwork/telescope/internal/converter/diag"
	"github.com/blockopsnetwork/telescope/internal/converter/internal/common"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"go.opentelemetry.io/collector/component"
)

func init() {
	converters = append(converters, transformProcessorConverter{})
}

type transformProcessorConverter struct{}

func (transformProcessorConverter) Factory() component.Factory {
	return transformprocessor.NewFactory()
}

func (transformProcessorConverter) InputComponentName() string {
	return "otelcol.processor.transform"
}

func (transformProcessorConverter) ConvertAndAppend(state *State, id component.InstanceID, cfg component.Config) diag.Diagnostics {
	var diags diag.Diagnostics

	label := state.FlowComponentLabel()

	args := toTransformProcessor(state, id, cfg.(*transformprocessor.Config))
	block := common.NewBlockWithOverride([]string{"otelcol", "processor", "transform"}, label, args)

	diags.Add(
		diag.SeverityLevelInfo,
		fmt.Sprintf("Converted %s into %s", StringifyInstanceID(id), StringifyBlock(block)),
	)

	state.Body().AppendBlock(block)
	return diags
}

func toTransformProcessor(state *State, id component.InstanceID, cfg *transformprocessor.Config) *transform.Arguments {
	var (
		nextMetrics = state.Next(id, component.DataTypeMetrics)
		nextLogs    = state.Next(id, component.DataTypeLogs)
		nextTraces  = state.Next(id, component.DataTypeTraces)
	)

	return &transform.Arguments{
		ErrorMode:        cfg.ErrorMode,
		TraceStatements:  toContextStatements(encodeMapslice(cfg.TraceStatements)),
		MetricStatements: toContextStatements(encodeMapslice(cfg.MetricStatements)),
		LogStatements:    toContextStatements(encodeMapslice(cfg.LogStatements)),
		Output: &otelcol.ConsumerArguments{
			Metrics: ToTokenizedConsumers(nextMetrics),
			Logs:    ToTokenizedConsumers(nextLogs),
			Traces:  ToTokenizedConsumers(nextTraces),
		},
	}
}

func toContextStatements(in []map[string]any) []transform.ContextStatements {
	res := make([]transform.ContextStatements, 0, len(in))
	for _, s := range in {
		res = append(res, transform.ContextStatements{
			Context:    transform.ContextID(encodeString(s["context"])),
			Statements: s["statements"].([]string),
		})
	}

	return res
}
