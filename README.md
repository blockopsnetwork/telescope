<p align="center">
  <a href="https://app.blockops.network" title="Blockops Network">
    <img src="./assets/img/blockops-logo.png" alt="Blockops-Network-logo" width="244" />
  </a>
</p>

<h1 align="center">Telescope</h1>

- [Summary](#summary)
- [Usage](#usage)
  - [Basic Configuration](#basic-configuration)
  - [Ethereum Integration](#ethereum-integration)
  - [Using Configuration File](#using-configuration-file)
- [Supported Networks](#supported-networks)
- [Language](#language)
- [Contributing](#contributing)
- [License](#license)


## Summary

An All-in-One Web3 Observability tooling that collects metrics and logs from blockchain nodes and related infrastructure.

## Usage

Telescope can be configured either through command line flags or a YAML configuration file.

### Basic Configuration

Basic usage with metrics enabled:

```bash
telescope \
  --metrics \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-name \
  --telescope-username=user \
  --telescope-password=pass \
  --remote-write-url=https://metrics.example.com
```

Enable both metrics and logs:

```bash
telescope \
  --metrics \
  --enable-logs \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-name \
  --telescope-username=user \
  --telescope-password=pass \
  --remote-write-url=https://metrics.example.com \
  --logs-sink-url=https://logs.example.com \
  --telescope-loki-username=loki-user \
  --telescope-loki-password=loki-pass
```

### Ethereum Integration

Telescope includes native Ethereum blockchain metrics collection that replaces the need for running a separate `ethereum-metrics-exporter`. This integration supports both execution and consensus layer monitoring.

**⚠️ Important**: Ethereum integration requires the `--enable-features integrations-next` flag.

#### Basic Ethereum Integration

Monitor Ethereum execution and consensus nodes:

```bash
telescope \
  --enable-features integrations-next \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-project \
  --telescope-username=user \
  --telescope-password=pass \
  --remote-write-url=https://prometheus.example.com/api/v1/write \
  --ethereum-execution-url=http://localhost:8545 \
  --ethereum-consensus-url=http://localhost:5052
```

#### Ethereum Integration with Custom Modules

Configure which execution client modules to monitor:

```bash
telescope \
  --enable-features integrations-next \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-project \
  --telescope-username=user \
  --telescope-password=pass \
  --remote-write-url=https://prometheus.example.com/api/v1/write \
  --ethereum-execution-url=http://localhost:8545 \
  --ethereum-execution-modules=sync,eth,net,web3,txpool \
  --ethereum-consensus-url=http://localhost:5052
```

#### Ethereum Integration with Disk Usage Monitoring

Enable separate disk usage monitoring (useful when node_exporter is not available):

```bash
telescope \
  --enable-features integrations-next \
  --network=ssv \
  --project-id=my-project \
  --project-name=my-project \
  --telescope-username=user \
  --telescope-password=pass \
  --remote-write-url=https://prometheus.example.com/api/v1/write \
  --ethereum-execution-url=http://localhost:8545 \
  --ethereum-disk-usage-enabled \
  --ethereum-disk-usage-dirs=/data/ethereum,/data/consensus \
  --ethereum-disk-usage-interval=10m
```

#### Available Ethereum Flags

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `--ethereum-enabled` | Enable Ethereum metrics collection | `false` | No |
| `--ethereum-execution-url` | Ethereum execution node URL | - | No¹ |
| `--ethereum-consensus-url` | Ethereum consensus node URL | - | No¹ |
| `--ethereum-execution-modules` | Execution modules to enable | `sync,eth,net,web3,txpool` | No |
| `--ethereum-disk-usage-enabled` | Enable disk usage monitoring | `false` | No |
| `--ethereum-disk-usage-dirs` | Directories to monitor | - | Yes² |
| `--ethereum-disk-usage-interval` | Disk usage collection interval | `5m` | No |

¹ At least one of `--ethereum-execution-url`, `--ethereum-consensus-url`, or `--ethereum-disk-usage-enabled` must be provided when using Ethereum integration.

² Required when `--ethereum-disk-usage-enabled` is true.

#### Collected Metrics

The Ethereum integration collects metrics with the `eth_exe_` prefix for execution layer and `eth_con_` prefix for consensus layer, including:

- **Execution Layer**: Block height, peer count, sync status, transaction pool metrics, and more
- **Consensus Layer**: Validator metrics, attestation performance, sync committee participation
- **Disk Usage**: Directory size monitoring with configurable intervals

### Using Configuration File

Create a YAML configuration file and run:

```bash
telescope --config-file=telescope_config.yaml
```

Example configuration file:

```yaml
server:
  log_level: info
metrics:
  global:
    scrape_interval: 15s
    external_labels:
      project_id: my-project
      project_name: my-name
    remote_write:
      - url: https://metrics.example.com
        basic_auth:
          username: user
          password: pass
  wal_directory: /tmp/telescope
  configs:
    - name: my-name_ethereum_metrics
      host_filter: false
      scrape_configs:
        - job_name: ethereum
          static_configs:
            - targets: ["localhost:8545"]
logs:
  configs:
    - name: telescope_logs
      clients:
        - url: https://logs.example.com
          basic_auth:
            username: loki-user
            password: loki-pass
          external_labels:
            project_id: my-project
            project_name: my-name
      positions:
        filename: /tmp/telescope_logs
integrations:
  agent:
    autoscrape:
      enable: true
      metrics_instance: "my-name_ethereum_metrics"
  node_exporter:
    autoscrape:
      enable: true
      metrics_instance: "my-name_ethereum_metrics"
  # Optional: Ethereum integration (requires --enable-features integrations-next)
  ethereum_configs:
    - instance: "ethereum_node_1"
      enabled: true
      autoscrape:
        enable: true
        metrics_instance: "my-name_ethereum_metrics"
      execution:
        enabled: true
        url: "http://localhost:8545"
        modules: ["sync", "eth", "net", "web3", "txpool"]
      consensus:
        enabled: true
        url: "http://localhost:5052"
        event_stream:
          enabled: true
          topics: ["head", "finalized_checkpoint"]
      disk_usage:
        enabled: true
        directories: ["/data/ethereum", "/data/consensus"]
        interval: "5m"
```

## Supported Networks

Telescope supports the following blockchain networks:

- **ethereum**: Ethereum mainnet and testnets
- **polkadot**: Polkadot ecosystem
- **hyperbridge**: Hyperbridge network
- **ssv**: Secret Shared Validators (SSV) network

Use the `--network` flag to specify which network configuration to use.



## Language
- Golang

## Contributing
We would love to work with anyone who can contribute their work and improve this project. The details will be shared soon.


## License

Licensed Under [Apache 2.0](./LICENSE)