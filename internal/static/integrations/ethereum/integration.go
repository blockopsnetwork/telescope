package ethereum

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethpandaops/beacon/pkg/beacon"
	"github.com/ethpandaops/ethereum-metrics-exporter/pkg/exporter/disk"
	"github.com/ethpandaops/ethereum-metrics-exporter/pkg/exporter/execution/api"
	"github.com/ethpandaops/ethereum-metrics-exporter/pkg/exporter/execution/jobs"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/onrik/ethrpc"
	v2integrations "github.com/blockopsnetwork/telescope/internal/static/integrations/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promConfig "github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/common/model"
	"github.com/blockopsnetwork/telescope/internal/static/integrations/v2/autoscrape"
	"github.com/sirupsen/logrus"
)

// Integration implements the ethereum integration
// Ensure it implements the required interfaces
var (
	_ v2integrations.Integration        = (*Integration)(nil)
	_ v2integrations.HTTPIntegration    = (*Integration)(nil)
	_ v2integrations.MetricsIntegration = (*Integration)(nil)
)

type Integration struct {
	log log.Logger
	cfg *Config
	reg prometheus.Registerer
	globals v2integrations.Globals

	// Clients
	ethClient    *ethclient.Client
	ethRPCClient *ethrpc.EthRPC
	beaconClient beacon.Node

	// Metrics collectors
	syncMetrics    jobs.SyncStatus
	generalMetrics jobs.GeneralMetrics
	txpoolMetrics  jobs.TXPool
	adminMetrics   jobs.Admin
	blockMetrics   jobs.BlockMetrics
	web3Metrics    jobs.Web3
	netMetrics     jobs.Net
	diskUsage      disk.UsageMetrics

	// Track if metrics are already registered
	metricsRegistered bool
}

// New creates a new ethereum integration
func New(log log.Logger, cfg *Config, globals v2integrations.Globals) *Integration {
	return &Integration{
		log:     log,
		cfg:     cfg,
		reg:     prometheus.DefaultRegisterer,
		globals: globals,
	}
}

// RunIntegration implements v2integrations.Integration
func (i *Integration) RunIntegration(ctx context.Context) error {
	if !i.cfg.Enabled {
		level.Info(i.log).Log("msg", "ethereum integration disabled")
		return nil
	}

	// Prevent double initialization
	if i.metricsRegistered {
		level.Info(i.log).Log("msg", "ethereum integration already running")
		<-ctx.Done()
		return nil
	}

	level.Info(i.log).Log("msg", "starting ethereum integration")

	if i.cfg.Execution.Enabled {
		if err := i.setupExecutionClient(ctx); err != nil {
			return fmt.Errorf("failed to setup execution client: %w", err)
		}
	}

	if i.cfg.Consensus.Enabled {
		if err := i.setupConsensusClient(ctx); err != nil {
			return fmt.Errorf("failed to setup consensus client: %w", err)
		}
		// Register consensus metrics if available
		if i.beaconClient != nil {
			if collector, ok := i.beaconClient.(prometheus.Collector); ok {
				if err := i.reg.Register(collector); err != nil {
					level.Error(i.log).Log("msg", "failed to register consensus metrics", "err", err)
				}
			}
		}
	}

	if i.cfg.DiskUsage.Enabled {
		if err := i.setupDiskUsage(ctx); err != nil {
			return fmt.Errorf("failed to setup disk usage: %w", err)
		}
		// Register disk usage metrics
		if i.diskUsage != nil {
			if collector, ok := i.diskUsage.(prometheus.Collector); ok {
				if err := i.reg.Register(collector); err != nil {
					level.Error(i.log).Log("msg", "failed to register disk usage metrics", "err", err)
				}
			}
		}
	}

	// Mark as registered to prevent duplicate setup
	i.metricsRegistered = true
	
	// Wait for context cancellation
	<-ctx.Done()
	return nil
}

func (i *Integration) setupExecutionClient(ctx context.Context) error {
	level.Info(i.log).Log("msg", "setting up execution client", "url", i.cfg.Execution.URL)

	// Create Ethereum client
	client, err := ethclient.Dial(i.cfg.Execution.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to execution client: %w", err)
	}
	i.ethClient = client

	// Create RPC client
	rpcClient := ethrpc.NewEthRPC(i.cfg.Execution.URL)
	i.ethRPCClient = rpcClient

	// Create internal API client
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(log.NewStdlibAdapter(i.log))
	internalAPI := api.NewExecutionClient(ctx, logrusLogger, i.cfg.Execution.URL)

	// Create const labels matching the original exporter exactly
	constLabels := make(prometheus.Labels)
	constLabels["ethereum_role"] = "execution"
	constLabels["node_name"] = "ethereum"

	// Initialize all metrics collectors exactly like the original exporter
	i.syncMetrics = jobs.NewSyncStatus(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth_exe", constLabels)
	i.generalMetrics = jobs.NewGeneralMetrics(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth_exe", constLabels)
	i.txpoolMetrics = jobs.NewTXPool(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth_exe", constLabels)
	i.adminMetrics = jobs.NewAdmin(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth_exe", constLabels)
	i.blockMetrics = jobs.NewBlockMetrics(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth_exe", constLabels)
	i.web3Metrics = jobs.NewWeb3(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth_exe", constLabels)
	i.netMetrics = jobs.NewNet(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth_exe", constLabels)

	// Enable and register metrics based on modules - exactly like the original
	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, i.syncMetrics.RequiredModules()); able {
		level.Info(i.log).Log("msg", "Enabling sync status metrics")
		i.reg.MustRegister(i.syncMetrics.Percentage)
		i.reg.MustRegister(i.syncMetrics.StartingBlock)
		i.reg.MustRegister(i.syncMetrics.CurrentBlock)
		i.reg.MustRegister(i.syncMetrics.IsSyncing)
		i.reg.MustRegister(i.syncMetrics.HighestBlock)
		go i.syncMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, i.generalMetrics.RequiredModules()); able {
		level.Info(i.log).Log("msg", "Enabling general metrics")
		i.reg.MustRegister(i.generalMetrics.NetworkID)
		i.reg.MustRegister(i.generalMetrics.GasPrice)
		i.reg.MustRegister(i.generalMetrics.ChainID)
		go i.generalMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, i.blockMetrics.RequiredModules()); able {
		level.Info(i.log).Log("msg", "Enabling block metrics")
		i.reg.MustRegister(i.blockMetrics.MostRecentBlockNumber)
		i.reg.MustRegister(i.blockMetrics.HeadBlockSize)
		i.reg.MustRegister(i.blockMetrics.HeadGasLimit)
		i.reg.MustRegister(i.blockMetrics.HeadGasUsed)
		i.reg.MustRegister(i.blockMetrics.HeadTransactionCount)
		i.reg.MustRegister(i.blockMetrics.HeadBaseFeePerGas)
		i.reg.MustRegister(i.blockMetrics.SafeBaseFeePerGas)
		i.reg.MustRegister(i.blockMetrics.SafeBlockSize)
		i.reg.MustRegister(i.blockMetrics.SafeGasLimit)
		i.reg.MustRegister(i.blockMetrics.SafeGasUsed)
		i.reg.MustRegister(i.blockMetrics.SafeTransactionCount)
		go i.blockMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, i.txpoolMetrics.RequiredModules()); able {
		level.Info(i.log).Log("msg", "Enabling txpool metrics")
		i.reg.MustRegister(i.txpoolMetrics.Transactions)
		go i.txpoolMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, i.adminMetrics.RequiredModules()); able {
		level.Info(i.log).Log("msg", "Enabling admin metrics")
		i.reg.MustRegister(i.adminMetrics.NodeInfo)
		i.reg.MustRegister(i.adminMetrics.Port)
		i.reg.MustRegister(i.adminMetrics.Peers)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					level.Error(i.log).Log("msg", "admin metrics crashed", "err", r)
				}
			}()
			i.adminMetrics.Start(ctx)
		}()
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, i.web3Metrics.RequiredModules()); able {
		level.Info(i.log).Log("msg", "Enabling web3 metrics")
		i.reg.MustRegister(i.web3Metrics.ClientVersion)
		go i.web3Metrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, i.netMetrics.RequiredModules()); able {
		level.Info(i.log).Log("msg", "Enabling net metrics")
		i.reg.MustRegister(i.netMetrics.PeerCount)
		go i.netMetrics.Start(ctx)
	}

	return nil
}

func (i *Integration) setupConsensusClient(ctx context.Context) error {
	level.Info(i.log).Log("msg", "setting up consensus client", "url", i.cfg.Consensus.URL)

	// Create a logrus logger from our go-kit logger
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(log.NewStdlibAdapter(i.log))

	opts := *beacon.DefaultOptions().EnablePrometheusMetrics()
	if i.cfg.Consensus.EventStream.Enabled {
		opts.BeaconSubscription.Topics = i.cfg.Consensus.EventStream.Topics
		if len(opts.BeaconSubscription.Topics) == 0 {
			opts.EnableDefaultBeaconSubscription()
		}
		opts.BeaconSubscription.Enabled = true
	}

	i.beaconClient = beacon.NewNode(logrusLogger, &beacon.Config{
		Addr: i.cfg.Consensus.URL,
		Name: "ethereum",
	}, "eth_con", opts)

	return nil
}

type diskUsageWrapper struct {
	disk.UsageMetrics
	metrics prometheus.Collector
}

func (d *diskUsageWrapper) Collect(ch chan<- prometheus.Metric) {
	if d.metrics != nil {
		d.metrics.Collect(ch)
	}
}

func (d *diskUsageWrapper) Describe(ch chan<- *prometheus.Desc) {
	if d.metrics != nil {
		d.metrics.Describe(ch)
	}
}

func (i *Integration) setupDiskUsage(ctx context.Context) error {
	level.Info(i.log).Log("msg", "setting up disk usage monitoring", "directories", i.cfg.DiskUsage.Directories)

	interval, err := time.ParseDuration(i.cfg.DiskUsage.Interval)
	if err != nil {
		return fmt.Errorf("invalid disk usage interval: %w", err)
	}

	// Create a logrus logger from our go-kit logger
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(log.NewStdlibAdapter(i.log))

	diskUsage, err := disk.NewUsage(
		ctx,
		logrusLogger,
		"eth_disk",
		i.cfg.DiskUsage.Directories,
		interval,
	)
	if err != nil {
		return err
	}

	// Get the metrics collector from the diskUsage instance, if possible
	var metricsCollector prometheus.Collector
	if mGetter, ok := diskUsage.(interface{ MetricsCollector() prometheus.Collector }); ok {
		metricsCollector = mGetter.MetricsCollector()
	}

	i.diskUsage = &diskUsageWrapper{
		UsageMetrics: diskUsage,
		metrics:      metricsCollector,
	}

	go diskUsage.StartAsync(ctx)

	return nil
}

// Stop stops the integration
func (i *Integration) Stop() {
	if i.ethClient != nil {
		i.ethClient.Close()
	}
	if i.beaconClient != nil {
		if stopper, ok := i.beaconClient.(interface{ Stop(context.Context) error }); ok {
			_ = stopper.Stop(context.Background())
		}
	}
	if i.diskUsage != nil {
		if stopper, ok := i.diskUsage.(interface{ Stop() }); ok {
			stopper.Stop()
		}
	}
}

// Handler implements v2integrations.HTTPIntegration
func (i *Integration) Handler(prefix string) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.Handle(prefix+"/metrics", promhttp.Handler())
	return mux, nil
}

// Targets implements v2integrations.MetricsIntegration
func (i *Integration) Targets(ep v2integrations.Endpoint) []*targetgroup.Group {
	return []*targetgroup.Group{
		{
			Targets: []model.LabelSet{
				{model.AddressLabel: model.LabelValue(ep.Host)},
			},
			Labels: model.LabelSet{
				"job":      "ethereum",
				"instance": model.LabelValue(ep.Host),
			},
		},
	}
}

// ScrapeConfigs implements v2integrations.MetricsIntegration
func (i *Integration) ScrapeConfigs(sd discovery.Configs) []*autoscrape.ScrapeConfig {
	// Check if autoscrape is enabled
	if !*i.cfg.Common.Autoscrape.Enable {
		return nil
	}

	// Create scrape config based on the default
	cfg := promConfig.DefaultScrapeConfig
	cfg.JobName = fmt.Sprintf("ethereum/%s", i.cfg.Name())
	cfg.Scheme = i.globals.AgentBaseURL.Scheme
	cfg.ServiceDiscoveryConfigs = sd
	cfg.ScrapeInterval = i.cfg.Common.Autoscrape.ScrapeInterval
	cfg.ScrapeTimeout = i.cfg.Common.Autoscrape.ScrapeTimeout
	cfg.RelabelConfigs = i.cfg.Common.Autoscrape.RelabelConfigs
	cfg.MetricRelabelConfigs = i.cfg.Common.Autoscrape.MetricRelabelConfigs

	return []*autoscrape.ScrapeConfig{{
		Instance: i.cfg.Common.Autoscrape.MetricsInstance,
		Config:   cfg,
	}}
}
