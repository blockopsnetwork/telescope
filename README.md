# Telescope

<p align="center">
  <a href="https://app.blockops.network" title="Blockops Network">
    <img src="./assets/img/blockops-logo.png" alt="Blockops-Network-logo" width="244" />
  </a>
</p>

- [Summary](#summary)
- [Usage](#usage)
  - [Basic Configuration](#basic-configuration)
  - [Auto-Discovery and Configuration Generation](#auto-discovery-and-configuration-generation)
  - [Logs Collection](#logs-collection)
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

Basic usage with metrics enabled and auto-discovery:

```bash
telescope \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-project \
  --telescope-username=user \
  --telescope-password=pass \
  --remote-write-url=https://prometheus.example.com/api/v1/write
```

Enable both metrics and logs with auto-discovery:

```bash
telescope \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-project \
  --telescope-username=user \
  --telescope-password=pass \
  --remote-write-url=https://prometheus.example.com/api/v1/write \
  --enable-logs=true \
  --logs-sink-url=https://loki.example.com/loki/api/v1/push \
  --telescope-loki-username=user \
  --telescope-loki-password=pass
```

For Docker environments with container log collection:

```bash
telescope \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-project \
  --telescope-username=user \
  --telescope-password=pass \
  --remote-write-url=https://prometheus.example.com/api/v1/write \
  --enable-logs=true \
  --enable-docker-logs=true \
  --logs-sink-url=https://loki.example.com/loki/api/v1/push \
  --telescope-loki-username=user \
  --telescope-loki-password=pass
```

### Auto-Discovery and Configuration Generation

Telescope features intelligent auto-discovery and configuration generation that automatically creates comprehensive monitoring configurations based on your network and requirements. Instead of manually writing complex YAML configurations, you can use command-line flags and Telescope will generate the complete configuration automatically.

#### How Auto-Discovery Works

1. **Network Detection**: Based on the `--network` flag, Telescope automatically discovers and configures appropriate scrape targets for your blockchain network
2. **Service Configuration**: Automatically configures metrics collection, log aggregation, and integrations based on enabled features
3. **Target Generation**: Creates scrape configs with proper job names, targets, and labeling for your infrastructure

#### Generated Configuration File

When you run Telescope with command-line flags, it automatically generates a `telescope_config.yaml` file containing:

- **Metrics Configuration**: Prometheus-compatible scrape configs with proper intervals and labeling
- **Logs Configuration**: Loki client configuration and log collection rules
- **Integration Configuration**: Enabled integrations (Node Exporter, Ethereum, etc.) with autoscrape
- **Network-Specific Targets**: Automatically discovered endpoints based on your network choice

Example of auto-generated configuration for Polkadot:

```bash
telescope --network=polkadot --project-id=test --project-name=test \
    --telescope-username=user --telescope-password=pass \
    --remote-write-url=https://example.com/write \
    --enable-features integrations-next
```

This generates a complete configuration with:
- Polkadot relay chain monitoring (port 30333)
- Parachain monitoring (port 9933)  
- Node exporter integration
- Proper labeling and external labels

#### Supported Auto-Discovery Networks

| Network | Targets Discovered | Default Ports |
|---------|-------------------|---------------|
| `ethereum` | Execution + Consensus nodes | 6060, 8008 |
| `polkadot` | Relay chain + Parachains | 30333, 9933 |
| `hyperbridge` | Hyperbridge node | 8080 |
| `ssv` | Execution + Consensus + MEV-Boost + SSV-DKG + SSV node | 6060, 8008, 18550, 3030, 13000 |

### Logs Collection

Telescope supports comprehensive log collection with automatic configuration generation. You can enable basic log collection or advanced Docker container log scraping.

#### Basic Log Collection

Enable basic log collection to send application logs to Loki:

```bash
telescope \
  --enable-logs=true \
  --logs-sink-url=https://loki.example.com/loki/api/v1/push \
  --telescope-loki-username=user \
  --telescope-loki-password=pass \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-project
```

#### Docker Container Log Scraping

For containerized environments, enable Docker log scraping to automatically collect logs from all Docker containers:

```bash
telescope \
  --enable-logs=true \
  --enable-docker-logs=true \
  --logs-sink-url=https://loki.example.com/loki/api/v1/push \
  --telescope-loki-username=user \
  --telescope-loki-password=pass \
  --docker-host=unix:///var/run/docker.sock \
  --network=ethereum \
  --project-id=my-project \
  --project-name=my-project
```

#### Generated Log Configuration

When Docker logs are enabled, Telescope automatically generates:

- **Docker Service Discovery**: Connects to Docker daemon for container discovery
- **Comprehensive Relabeling**: Extracts container metadata as labels
- **Label Mapping**: Maps Docker labels to log labels for filtering and organization

The generated configuration includes relabel rules for:
- Container name and log stream
- Custom Docker labels (network, client_name, group, host_type, etc.)
- Project identification labels
- Instance and location labels

#### Available Log Flags

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `--enable-logs` | Enable log collection | `false` | No |
| `--logs-sink-url` | Loki endpoint URL | - | Yes¹ |
| `--telescope-loki-username` | Loki authentication username | - | No |
| `--telescope-loki-password` | Loki authentication password | - | No |
| `--enable-docker-logs` | Enable Docker container log scraping | `false` | No |
| `--docker-host` | Docker daemon socket | `unix:///var/run/docker.sock` | No |

¹ Required when `--enable-logs=true`

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


#### Available Ethereum Flags

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `--ethereum-enabled` | Enable Ethereum metrics collection | `false` | No |
| `--ethereum-execution-url` | Ethereum execution node URL | - | No¹ |
| `--ethereum-consensus-url` | Ethereum consensus node URL | - | No¹ |
| `--ethereum-execution-modules` | Execution modules to enable | `sync,eth,net,web3,txpool` | No |

¹ At least one of `--ethereum-execution-url` or `--ethereum-consensus-url` must be provided when using Ethereum integration.

#### Collected Metrics

The Ethereum integration collects metrics with the `eth_exe_` prefix for execution layer and `eth_con_` prefix for consensus layer, including:

- **Execution Layer**: Block height, peer count, sync status, transaction pool metrics, and more
- **Consensus Layer**: Validator metrics, attestation performance, sync committee participation

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