package secret

import (
	"github.com/blockopsnetwork/telescope/internal/component"
	"github.com/blockopsnetwork/telescope/internal/component/remote/kubernetes"
	"github.com/blockopsnetwork/telescope/internal/featuregate"
)

func init() {
	component.Register(component.Registration{
		Name:      "remote.kubernetes.secret",
		Stability: featuregate.StabilityStable,
		Args:      kubernetes.Arguments{},
		Exports:   kubernetes.Exports{},
		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			return kubernetes.New(opts, args.(kubernetes.Arguments), kubernetes.TypeSecret)
		},
	})
}
