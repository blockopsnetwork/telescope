package probes

import (
	"github.com/blockopsnetwork/telescope/internal/component"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/operator"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/operator/common"
	"github.com/blockopsnetwork/telescope/internal/featuregate"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.operator.probes",
		Stability: featuregate.StabilityBeta,
		Args:      operator.Arguments{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			return common.New(opts, args, common.KindProbe)
		},
	})
}
