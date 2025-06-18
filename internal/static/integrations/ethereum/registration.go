package ethereum

import (
	"github.com/blockopsnetwork/telescope/internal/static/integrations"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	integrations.Register("ethereum", func(log log.Logger, reg prometheus.Registerer, cfg interface{}) (integrations.Integration, error) {
		config, ok := cfg.(*Config)
		if !ok {
			return nil, integrations.ErrInvalidConfig
		}

		// Create a new registry that wraps the provided registerer
		registry := prometheus.NewRegistry()
		if reg != nil {
			// If a registerer was provided, use it
			registry = reg.(*prometheus.Registry)
		}

		return New(log, config, registry), nil
	})
}
