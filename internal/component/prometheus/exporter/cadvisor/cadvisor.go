package cadvisor

import (
	"time"

	"github.com/blockopsnetwork/telescope/internal/component"
	"github.com/blockopsnetwork/telescope/internal/component/prometheus/exporter"
	"github.com/blockopsnetwork/telescope/internal/featuregate"
	"github.com/blockopsnetwork/telescope/internal/static/integrations"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/cadvisor"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.exporter.cadvisor",
		Stability: featuregate.StabilityStable,
		Args:      Arguments{},
		Exports:   exporter.Exports{},

		Build: exporter.New(createExporter, "cadvisor"),
	})
}

func createExporter(opts component.Options, args component.Arguments, defaultInstanceKey string) (integrations.Integration, string, error) {
	a := args.(Arguments)
	return integrations.NewIntegrationWithInstanceKey(opts.Logger, a.Convert(), defaultInstanceKey)
}

// Arguments configures the prometheus.exporter.cadvisor component.
type Arguments struct {
	StoreContainerLabels       bool          `river:"store_container_labels,attr,optional"`
	AllowlistedContainerLabels []string      `river:"allowlisted_container_labels,attr,optional"`
	EnvMetadataAllowlist       []string      `river:"env_metadata_allowlist,attr,optional"`
	RawCgroupPrefixAllowlist   []string      `river:"raw_cgroup_prefix_allowlist,attr,optional"`
	PerfEventsConfig           string        `river:"perf_events_config,attr,optional"`
	ResctrlInterval            time.Duration `river:"resctrl_interval,attr,optional"`
	DisabledMetrics            []string      `river:"disabled_metrics,attr,optional"`
	EnabledMetrics             []string      `river:"enabled_metrics,attr,optional"`
	StorageDuration            time.Duration `river:"storage_duration,attr,optional"`
	ContainerdHost             string        `river:"containerd_host,attr,optional"`
	ContainerdNamespace        string        `river:"containerd_namespace,attr,optional"`
	DockerHost                 string        `river:"docker_host,attr,optional"`
	UseDockerTLS               bool          `river:"use_docker_tls,attr,optional"`
	DockerTLSCert              string        `river:"docker_tls_cert,attr,optional"`
	DockerTLSKey               string        `river:"docker_tls_key,attr,optional"`
	DockerTLSCA                string        `river:"docker_tls_ca,attr,optional"`
	DockerOnly                 bool          `river:"docker_only,attr,optional"`
	DisableRootCgroupStats     bool          `river:"disable_root_cgroup_stats,attr,optional"`
}

// SetToDefault implements river.Defaulter.
func (a *Arguments) SetToDefault() {
	*a = Arguments{
		StoreContainerLabels:       true,
		AllowlistedContainerLabels: []string{""},
		EnvMetadataAllowlist:       []string{""},
		RawCgroupPrefixAllowlist:   []string{""},
		ResctrlInterval:            0,
		StorageDuration:            2 * time.Minute,

		ContainerdHost:      "/run/containerd/containerd.sock",
		ContainerdNamespace: "k8s.io",

		// TODO(@tpaschalis) Do we need the default cert/key/ca since tls is disabled by default?
		DockerHost:    "unix:///var/run/docker.sock",
		UseDockerTLS:  false,
		DockerTLSCert: "cert.pem",
		DockerTLSKey:  "key.pem",
		DockerTLSCA:   "ca.pem",

		DockerOnly:             false,
		DisableRootCgroupStats: false,
	}
}

// Convert returns the upstream-compatible configuration struct.
func (a *Arguments) Convert() *cadvisor.Config {
	if len(a.AllowlistedContainerLabels) == 0 {
		a.AllowlistedContainerLabels = []string{""}
	}
	if len(a.RawCgroupPrefixAllowlist) == 0 {
		a.RawCgroupPrefixAllowlist = []string{""}
	}
	if len(a.EnvMetadataAllowlist) == 0 {
		a.EnvMetadataAllowlist = []string{""}
	}

	cfg := &cadvisor.Config{
		StoreContainerLabels:       a.StoreContainerLabels,
		AllowlistedContainerLabels: a.AllowlistedContainerLabels,
		EnvMetadataAllowlist:       a.EnvMetadataAllowlist,
		RawCgroupPrefixAllowlist:   a.RawCgroupPrefixAllowlist,
		PerfEventsConfig:           a.PerfEventsConfig,
		ResctrlInterval:            int64(a.ResctrlInterval), // TODO(@tpaschalis) This is so that the cadvisor package can re-cast back to time.Duration. Can we make it use time.Duration directly instead?
		DisabledMetrics:            a.DisabledMetrics,
		EnabledMetrics:             a.EnabledMetrics,
		StorageDuration:            a.StorageDuration,
		Containerd:                 a.ContainerdHost,
		ContainerdNamespace:        a.ContainerdNamespace,
		Docker:                     a.DockerHost,
		DockerTLS:                  a.UseDockerTLS,
		DockerTLSCert:              a.DockerTLSCert,
		DockerTLSKey:               a.DockerTLSKey,
		DockerTLSCA:                a.DockerTLSCA,
		DockerOnly:                 a.DockerOnly,
		DisableRootCgroupStats:     a.DisableRootCgroupStats,
	}

	return cfg
}
