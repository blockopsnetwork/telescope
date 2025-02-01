package component

import (
	"time"

	"github.com/blockopsnetwork/telescope/internal/component/discovery"
	"github.com/blockopsnetwork/telescope/internal/component/discovery/consul"
	"github.com/blockopsnetwork/telescope/internal/converter/diag"
	"github.com/blockopsnetwork/telescope/internal/converter/internal/common"
	"github.com/blockopsnetwork/telescope/internal/converter/internal/prometheusconvert/build"
	"github.com/grafana/river/rivertypes"
	prom_consul "github.com/prometheus/prometheus/discovery/consul"
)

func appendDiscoveryConsul(pb *build.PrometheusBlocks, label string, sdConfig *prom_consul.SDConfig) discovery.Exports {
	discoveryConsulArgs := toDiscoveryConsul(sdConfig)
	name := []string{"discovery", "consul"}
	block := common.NewBlockWithOverride(name, label, discoveryConsulArgs)
	pb.DiscoveryBlocks = append(pb.DiscoveryBlocks, build.NewPrometheusBlock(block, name, label, "", ""))
	return common.NewDiscoveryExports("discovery.consul." + label + ".targets")
}

func ValidateDiscoveryConsul(sdConfig *prom_consul.SDConfig) diag.Diagnostics {
	return common.ValidateHttpClientConfig(&sdConfig.HTTPClientConfig)
}

func toDiscoveryConsul(sdConfig *prom_consul.SDConfig) *consul.Arguments {
	if sdConfig == nil {
		return nil
	}

	return &consul.Arguments{
		Server:           sdConfig.Server,
		Token:            rivertypes.Secret(sdConfig.Token),
		Datacenter:       sdConfig.Datacenter,
		Namespace:        sdConfig.Namespace,
		Partition:        sdConfig.Partition,
		TagSeparator:     sdConfig.TagSeparator,
		Scheme:           sdConfig.Scheme,
		Username:         sdConfig.Username,
		Password:         rivertypes.Secret(sdConfig.Password),
		AllowStale:       sdConfig.AllowStale,
		Services:         sdConfig.Services,
		ServiceTags:      sdConfig.ServiceTags,
		NodeMeta:         sdConfig.NodeMeta,
		RefreshInterval:  time.Duration(sdConfig.RefreshInterval),
		HTTPClientConfig: *common.ToHttpClientConfig(&sdConfig.HTTPClientConfig),
	}
}
