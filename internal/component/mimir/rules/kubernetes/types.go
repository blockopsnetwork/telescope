package rules

import (
	"fmt"
	"time"

	"github.com/blockopsnetwork/telescope/internal/component/common/config"
	"github.com/blockopsnetwork/telescope/internal/component/common/kubernetes"
)

type Arguments struct {
	Address              string                  `river:"address,attr"`
	TenantID             string                  `river:"tenant_id,attr,optional"`
	UseLegacyRoutes      bool                    `river:"use_legacy_routes,attr,optional"`
	PrometheusHTTPPrefix string                  `river:"prometheus_http_prefix,attr,optional"`
	HTTPClientConfig     config.HTTPClientConfig `river:",squash"`
	SyncInterval         time.Duration           `river:"sync_interval,attr,optional"`
	MimirNameSpacePrefix string                  `river:"mimir_namespace_prefix,attr,optional"`

	RuleSelector          kubernetes.LabelSelector `river:"rule_selector,block,optional"`
	RuleNamespaceSelector kubernetes.LabelSelector `river:"rule_namespace_selector,block,optional"`
}

var DefaultArguments = Arguments{
	SyncInterval:         5 * time.Minute,
	MimirNameSpacePrefix: "agent",
	HTTPClientConfig:     config.DefaultHTTPClientConfig,
	PrometheusHTTPPrefix: "/prometheus",
}

// SetToDefault implements river.Defaulter.
func (args *Arguments) SetToDefault() {
	*args = DefaultArguments
}

// Validate implements river.Validator.
func (args *Arguments) Validate() error {
	if args.SyncInterval <= 0 {
		return fmt.Errorf("sync_interval must be greater than 0")
	}
	if args.MimirNameSpacePrefix == "" {
		return fmt.Errorf("mimir_namespace_prefix must not be empty")
	}

	// We must explicitly Validate because HTTPClientConfig is squashed and it won't run otherwise
	return args.HTTPClientConfig.Validate()
}
