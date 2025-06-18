package ethereum

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethpandaops/beacon/pkg/beacon"
	"github.com/ethpandaops/ethereum-metrics-exporter/pkg/exporter/disk"
	"github.com/ethpandaops/ethereum-metrics-exporter/pkg/exporter/execution"
	"github.com/ethpandaops/ethereum-metrics-exporter/pkg/exporter/execution/jobs"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/onrik/ethrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promConfig "github.com/prometheus/prometheus/config"
	"github.com/sirupsen/logrus"
)

// Integration implements the ethereum integration
type Integration struct {
	log log.Logger
	cfg *Config
	reg prometheus.Registerer

	// Clients
	ethClient    *ethclient.Client
	ethRPCClient *ethrpc.EthRPC
	beaconClient beacon.Node

	// Metrics collectors
	syncMetrics    *jobs.SyncStatus
	generalMetrics *jobs.GeneralMetrics
	txpoolMetrics  *jobs.TXPool
	adminMetrics   *jobs.Admin
	blockMetrics   *jobs.BlockMetrics
	web3Metrics    *jobs.Web3
	netMetrics     *jobs.Net
	diskUsage      disk.UsageMetrics

	// Add registry for metrics handler
	metricsRegistry *prometheus.Registry
}

// New creates a new ethereum integration
func New(log log.Logger, cfg *Config, reg prometheus.Registerer) *Integration {
	// Create a new registry that wraps the provided registerer
	registry := prometheus.NewRegistry()
	if reg != nil {
		// If a registerer was provided, use it
		registry = reg.(*prometheus.Registry)
	}

	return &Integration{
		log:             log,
		cfg:             cfg,
		reg:             reg,
		metricsRegistry: registry,
	}
}

// Run starts the integration
func (i *Integration) Run(ctx context.Context) error {
	if !i.cfg.Enabled {
		level.Info(i.log).Log("msg", "ethereum integration disabled")
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
				if err := i.metricsRegistry.Register(collector); err != nil {
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
			if err := i.metricsRegistry.Register(i.diskUsage); err != nil {
				level.Error(i.log).Log("msg", "failed to register disk usage metrics", "err", err)
			}
		}
	}

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
	internalAPI := execution.NewClient(i.ethClient, i.ethRPCClient)

	// Create const labels
	constLabels := make(prometheus.Labels)
	constLabels["ethereum_role"] = "execution"
	constLabels["node_name"] = "ethereum"

	// Create a logrus logger from our go-kit logger
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(log.NewStdlibAdapter(i.log))

	// Create and register metrics collectors
	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, []string{"eth"}); able {
		i.syncMetrics = jobs.NewSyncStatus(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth", constLabels)
		i.metricsRegistry.MustRegister(
			i.syncMetrics.Percentage,
			i.syncMetrics.StartingBlock,
			i.syncMetrics.CurrentBlock,
			i.syncMetrics.IsSyncing,
			i.syncMetrics.HighestBlock,
		)
		go i.syncMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, []string{"eth", "net"}); able {
		i.generalMetrics = jobs.NewGeneralMetrics(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth", constLabels)
		i.metricsRegistry.MustRegister(
			i.generalMetrics.NetworkID,
			i.generalMetrics.GasPrice,
			i.generalMetrics.ChainID,
		)
		go i.generalMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, []string{"txpool"}); able {
		i.txpoolMetrics = jobs.NewTXPool(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth", constLabels)
		i.metricsRegistry.MustRegister(i.txpoolMetrics.Transactions)
		go i.txpoolMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, []string{"admin"}); able {
		i.adminMetrics = jobs.NewAdmin(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth", constLabels)
		i.metricsRegistry.MustRegister(
			i.adminMetrics.NodeInfo,
			i.adminMetrics.Port,
			i.adminMetrics.Peers,
		)
		go i.adminMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, []string{"eth", "net"}); able {
		i.blockMetrics = jobs.NewBlockMetrics(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth", constLabels)
		i.metricsRegistry.MustRegister(
			i.blockMetrics.MostRecentBlockNumber,
			i.blockMetrics.HeadBlockSize,
			i.blockMetrics.HeadGasLimit,
			i.blockMetrics.HeadGasUsed,
			i.blockMetrics.HeadTransactionCount,
			i.blockMetrics.HeadBaseFeePerGas,
			i.blockMetrics.SafeBaseFeePerGas,
			i.blockMetrics.SafeBlockSize,
			i.blockMetrics.SafeGasLimit,
			i.blockMetrics.SafeGasUsed,
			i.blockMetrics.SafeTransactionCount,
		)
		go i.blockMetrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, []string{"web3"}); able {
		i.web3Metrics = jobs.NewWeb3(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth", constLabels)
		i.metricsRegistry.MustRegister(i.web3Metrics.ClientVersion)
		go i.web3Metrics.Start(ctx)
	}

	if able := jobs.ExporterCanRun(i.cfg.Execution.Modules, []string{"net"}); able {
		i.netMetrics = jobs.NewNet(i.ethClient, internalAPI, i.ethRPCClient, logrusLogger, "eth", constLabels)
		i.metricsRegistry.MustRegister(i.netMetrics.PeerCount)
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
	i.beaconClient = beacon.NewNode(logrusLogger, &beacon.Config{
		Addr: i.cfg.Consensus.URL,
		Name: "ethereum",
	}, "eth_con", opts)

	return nil
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

	i.diskUsage = diskUsage
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

// MetricsHandler returns an HTTP handler for the metrics endpoint
func (i *Integration) MetricsHandler() (http.Handler, error) {
	return promhttp.HandlerFor(i.metricsRegistry, promhttp.HandlerOpts{}), nil
}

// ScrapeConfigs tells Telescope how to scrape this integration's metrics
func (i *Integration) ScrapeConfigs() []promConfig.ScrapeConfig {
	return []promConfig.ScrapeConfig{{
		JobName:     "ethereum",
		MetricsPath: "/metrics",
		// Additional config (static_configs, relabel_configs, etc.) can be added here if needed
	}}
}
