module github.com/blockopsnetwork/telescope

go 1.23.0

retract (
	v1.3.191 // Published accidentally
	v1.2.99 // Published accidentally
	v1.2.99-rc1 // Published accidentally
)

require (
	cloud.google.com/go/pubsub v1.34.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.0-beta.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.4.0
	github.com/Azure/go-autorest/autorest v0.11.29
	github.com/IBM/sarama v1.45.2
	github.com/Lusitaniae/apache_exporter v0.11.1-0.20220518131644-f9522724dab4
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/PuerkitoBio/rehttp v1.3.0
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137
	github.com/aws/aws-sdk-go v1.50.27 // indirect
	github.com/aws/aws-sdk-go-v2 v1.36.5
	github.com/aws/aws-sdk-go-v2/config v1.27.4
	github.com/aws/aws-sdk-go-v2/service/s3 v1.49.0
	github.com/bmatcuk/doublestar v1.3.4
	github.com/burningalchemist/sql_exporter v0.0.0-20240103092044-466b38b6abc4
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/cilium/ebpf v0.12.3 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/coreos/go-systemd/v22 v22.5.0
	github.com/davidmparrott/kafka_exporter/v2 v2.0.1
	github.com/docker/docker v24.0.9+incompatible
	github.com/docker/go-connections v0.5.0
	github.com/drone/envsubst/v2 v2.0.0-20210730161058-179042472c46
	github.com/fatih/color v1.18.0
	github.com/fatih/structs v1.1.0
	github.com/fortytw2/leaktest v1.3.0
	github.com/fsnotify/fsnotify v1.7.0
	github.com/github/smimesign v0.2.0
	github.com/go-git/go-git/v5 v5.11.0
	github.com/go-kit/log v0.2.1
	github.com/go-logfmt/logfmt v0.6.0
	github.com/go-logr/logr v1.4.2
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.4
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb
	github.com/google/cadvisor v0.47.0
	github.com/google/dnsmasq_exporter v0.2.1-0.20230620100026-44b14480804a
	github.com/google/go-cmp v0.6.0
	github.com/google/go-jsonnet v0.18.0
	github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8
	github.com/google/renameio/v2 v2.0.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/grafana/ckit v0.0.0-20230906125525-c046c99a5c04
	github.com/grafana/cloudflare-go v0.0.0-20230110200409-c627cf6792f2
	github.com/grafana/dskit v0.0.0-20240104111617-ea101a3b86eb
	github.com/grafana/go-gelf/v2 v2.0.1
	github.com/grafana/loki v1.6.2-0.20240221085104-f9d188620153 // k190 branch
	github.com/grafana/pyroscope-go/godeltaprof v0.1.8
	github.com/grafana/pyroscope/api v0.4.0
	github.com/grafana/pyroscope/ebpf v0.4.3
	github.com/grafana/regexp v0.0.0-20221123153739-15dc172cd2db
	github.com/grafana/river v0.3.1-0.20240123144725-960753160cd1
	github.com/grafana/snowflake-prometheus-exporter v0.0.0-20221213150626-862cad8e9538
	github.com/grafana/tail v0.0.0-20230510142333-77b18831edf0
	github.com/grafana/vmware_exporter v0.0.5-beta
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/hashicorp/consul/api v1.27.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-discover v0.0.0-20230724184603-e89ebd1b2f65
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/golang-lru v1.0.2
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/hashicorp/vault/api v1.20.0
	github.com/hashicorp/vault/api/auth/approle v0.2.0
	github.com/hashicorp/vault/api/auth/aws v0.2.0
	github.com/hashicorp/vault/api/auth/azure v0.2.0
	github.com/hashicorp/vault/api/auth/gcp v0.5.0
	github.com/hashicorp/vault/api/auth/kubernetes v0.2.0
	github.com/hashicorp/vault/api/auth/ldap v0.10.0
	github.com/hashicorp/vault/api/auth/userpass v0.2.0
	github.com/heroku/x v0.0.61
	github.com/iamseth/oracledb_exporter v0.0.0-20230918193147-95e16f21ceee
	github.com/influxdata/go-syslog/v3 v3.0.1-0.20230911200830-875f5bc594a4
	github.com/jaegertracing/jaeger v1.54.0
	github.com/jmespath/go-jmespath v0.4.0
	github.com/json-iterator/go v1.1.12
	github.com/klauspost/compress v1.18.0
	github.com/lib/pq v1.10.9
	github.com/mackerelio/go-osstat v0.2.3
	github.com/miekg/dns v1.1.56
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c
	github.com/mitchellh/reflectwalk v1.0.2
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f
	github.com/ncabatoff/process-exporter v0.7.10
	github.com/nerdswords/yet-another-cloudwatch-exporter v0.55.0
	github.com/ohler55/ojg v1.20.1 // indirect
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/oliver006/redis_exporter v1.54.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/servicegraphconnector v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/basicauthextension v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/headerssetterextension v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/jaegerremotesampling v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oauth2clientauthextension v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/sigv4authextension v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/loki v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.95.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.96.0
	github.com/opentracing-contrib/go-grpc v0.0.0-20210225150812-73cb765af46e
	github.com/opentracing-contrib/go-stdlib v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/ory/dockertest/v3 v3.8.1
	github.com/oschwald/geoip2-golang v1.9.0
	github.com/percona/mongodb_exporter v0.39.1-0.20230706092307-28432707eb65
	github.com/phayes/freeport v0.0.0-20220201140144-74d24b5ae9f5
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2
	github.com/prometheus-community/elasticsearch_exporter v1.5.0
	github.com/prometheus-community/postgres_exporter v0.11.1
	github.com/prometheus-community/stackdriver_exporter v0.13.0
	github.com/prometheus-community/windows_exporter v0.24.1-0.20231127180936-5a872a227c2f
	github.com/prometheus-operator/prometheus-operator v0.66.0
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.66.0
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.66.0
	github.com/prometheus/blackbox_exporter v0.24.1-0.20230623125439-bd22efa1c900
	github.com/prometheus/client_golang v1.19.0
	github.com/prometheus/client_model v0.6.1
	github.com/prometheus/common v0.48.0
	github.com/prometheus/consul_exporter v0.8.0
	github.com/prometheus/memcached_exporter v0.13.0
	github.com/prometheus/mysqld_exporter v0.14.0
	github.com/prometheus/node_exporter v1.6.0
	github.com/prometheus/procfs v0.12.0
	github.com/prometheus/prometheus v1.99.0
	github.com/prometheus/snmp_exporter v0.24.1
	github.com/prometheus/statsd_exporter v0.22.8
	github.com/richardartoul/molecule v1.0.1-0.20221107223329-32cfee06a052
	github.com/rs/cors v1.10.1
	github.com/shirou/gopsutil/v3 v3.24.1
	github.com/sijms/go-ora/v2 v2.7.6
	github.com/sirupsen/logrus v1.9.3
	github.com/spaolacci/murmur3 v1.1.0
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.10.0
	github.com/testcontainers/testcontainers-go v0.27.0
	github.com/testcontainers/testcontainers-go/modules/k3s v0.0.0-20230615142642-c175df34bd1d
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/vincent-petithory/dataurl v1.0.0
	github.com/webdevops/azure-metrics-exporter v0.0.0-20230717202958-8701afc2b013
	github.com/webdevops/go-common v0.0.0-20231022162947-a6adfb05a7e9
	github.com/wk8/go-ordered-map v0.2.0
	github.com/xdg-go/scram v1.1.2
	github.com/zeebo/xxh3 v1.0.2
	go.opentelemetry.io/collector v0.96.0
	go.opentelemetry.io/collector/component v0.96.0
	go.opentelemetry.io/collector/config/configauth v0.96.0
	go.opentelemetry.io/collector/config/configcompression v0.96.0
	go.opentelemetry.io/collector/config/configgrpc v0.96.0
	go.opentelemetry.io/collector/config/confighttp v0.96.0
	go.opentelemetry.io/collector/config/confignet v0.96.0
	go.opentelemetry.io/collector/config/configopaque v1.3.0
	go.opentelemetry.io/collector/config/configtelemetry v0.96.0
	go.opentelemetry.io/collector/config/configtls v0.96.0
	go.opentelemetry.io/collector/confmap v0.96.0
	go.opentelemetry.io/collector/connector v0.96.0
	go.opentelemetry.io/collector/consumer v0.96.0
	go.opentelemetry.io/collector/exporter v0.96.0
	go.opentelemetry.io/collector/exporter/loggingexporter v0.96.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.96.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.96.0
	go.opentelemetry.io/collector/extension v0.96.0
	go.opentelemetry.io/collector/extension/auth v0.96.0
	go.opentelemetry.io/collector/featuregate v1.3.0
	go.opentelemetry.io/collector/otelcol v0.96.0
	go.opentelemetry.io/collector/pdata v1.34.0
	go.opentelemetry.io/collector/processor v0.96.0
	go.opentelemetry.io/collector/processor/batchprocessor v0.96.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.96.0
	go.opentelemetry.io/collector/receiver v0.96.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.96.0
	go.opentelemetry.io/collector/semconv v0.96.0
	go.opentelemetry.io/collector/service v0.96.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.45.0
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.24.0
	go.opentelemetry.io/otel/exporters/prometheus v0.46.0
	go.opentelemetry.io/otel/metric v1.34.0
	go.opentelemetry.io/otel/sdk v1.34.0
	go.opentelemetry.io/otel/sdk/metric v1.34.0
	go.opentelemetry.io/otel/trace v1.34.0
	go.opentelemetry.io/proto/otlp v1.1.0
	go.uber.org/atomic v1.11.0
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.38.0
	golang.org/x/exp v0.0.0-20240103183307-be819d1f06fc
	golang.org/x/net v0.40.0
	golang.org/x/oauth2 v0.26.0
	golang.org/x/sys v0.33.0
	golang.org/x/text v0.25.0
	golang.org/x/time v0.12.0
	google.golang.org/api v0.155.0
	google.golang.org/grpc v1.72.2
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.29.2
	k8s.io/apiextensions-apiserver v0.29.0
	k8s.io/client-go v0.29.2
	k8s.io/component-base v0.29.0
	k8s.io/klog/v2 v2.110.1
	k8s.io/utils v0.0.0-20240102154912-e7106e64919e
	sigs.k8s.io/controller-runtime v0.17.2
	sigs.k8s.io/yaml v1.4.0
)

require (
	cloud.google.com/go v0.112.0 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/iam v1.1.5 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/AlekSi/pointer v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go v66.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor v0.10.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph v0.8.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.2.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.23 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.12 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.6 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.0 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/ClickHouse/clickhouse-go v1.5.4 // indirect
	github.com/GehirnInc/crypt v0.0.0-20200316065508-bb7000b8a962 // indirect
	github.com/JohnCGriffin/overflow v0.0.0-20211019200055-46fa312c352c // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Microsoft/hcsshim v0.11.4 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230828082145-3c4c8a2d2371 // indirect
	github.com/alecthomas/participle/v2 v2.1.1 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/apache/arrow/go/v12 v12.0.1 // indirect
	github.com/apache/thrift v0.19.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/avvmoto/buf-readerat v0.0.0-20171115124131-a17c8cb89270 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.32
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.1 // indirect
	github.com/aws/smithy-go v1.22.4 // indirect
	github.com/beevik/ntp v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.2-0.20180723201105-3c1074078d32+incompatible // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/boynux/squid-exporter v1.10.5-0.20230618153315-c1fae094e18e
	github.com/c2h5oh/datasize v0.0.0-20220606134207-859f65c6625b // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/checkpoint-restore/go-criu/v5 v5.3.0 // indirect
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 // indirect
	github.com/cncf/xds/go v0.0.0-20250121191232-2f005788dc42 // indirect
	github.com/containerd/cgroups v1.1.0 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/containerd/containerd v1.7.11 // indirect
	github.com/containerd/continuity v0.4.2 // indirect
	github.com/containerd/ttrpc v1.2.2 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/cpuguy83/dockercfg v0.3.1 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/danieljoos/wincred v1.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/dennwc/btrfs v0.0.0-20230312211831-a1f570bd01a1 // indirect
	github.com/dennwc/ioctl v1.0.0 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/denverdino/aliyungo v0.0.0-20190125010748-a747050bb1ba // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/digitalocean/godo v1.104.1 // indirect
	github.com/dimchansky/utfbom v1.1.1
	github.com/docker/cli v23.0.3+incompatible // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.6.0 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/efficientgo/tools/core v0.0.0-20220817170617-6c25e3b627dd // indirect
	github.com/elastic/go-sysinfo v1.8.1 // indirect
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/ema/qdisc v1.0.0 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.8.0 // indirect
	github.com/facette/natsort v0.0.0-20181210072756-2cd4dd1e2dcb // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.4 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/runtime v0.26.0 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/strfmt v0.21.7 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-openapi/validate v0.22.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/status v1.1.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gophercloud/gophercloud v1.7.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gosnmp/gosnmp v1.36.0 // indirect
	github.com/grafana/gomemcache v0.0.0-20231204155601-7de47a8c3cb0 // indirect
	github.com/grafana/loki/pkg/push v0.0.0-20231212100434-384e5c2dc872 // k180 branch
	github.com/grobie/gomemcache v0.0.0-20230213081705-239240bbc445 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/cronexpr v1.1.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-envparse v0.1.0 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/awsutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-7 // indirect
	github.com/hashicorp/mdns v1.0.4 // indirect
	github.com/hashicorp/memberlist v0.5.0 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20230721134942-515895c7690c // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hashicorp/vic v1.5.1-0.20190403131502-bbfe86ec9443 // indirect
	github.com/hodgesds/perf-utils v0.7.0 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/illumos/go-kstat v0.0.0-20210513183136-173c9b0a9973 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/infinityworks/go-common v0.0.0-20170820165359-7f20a140fd37 // indirect
	github.com/influxdata/telegraf v1.16.3 // indirect
	github.com/ionos-cloud/sdk-go/v6 v6.1.9 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.2 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/joyent/triton-go v0.0.0-20180628001255-830d2b111e62 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/jsimonetti/rtnetlink v1.3.5 // indirect
	github.com/karrick/godirwalk v1.17.0 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/krallistic/kazoo-go v0.0.0-20170526135507-a15279744f4e // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20181214104525-299bdde78165 // indirect
	github.com/linode/linodego v1.23.0 // indirect
	github.com/lufia/iostat v1.2.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mattn/go-xmlrpc v0.0.3 // indirect
	github.com/mdlayher/ethtool v0.1.0 // indirect
	github.com/mdlayher/genetlink v1.3.2 // indirect
	github.com/mdlayher/netlink v1.7.2 // indirect
	github.com/mdlayher/socket v0.4.1 // indirect
	github.com/mdlayher/wifi v0.1.0 // indirect
	github.com/microsoft/go-mssqldb v1.6.0 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mna/redisc v1.3.2 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mostynb/go-grpc-compression v1.2.2 // indirect
	github.com/mrunalp/fileutils v0.5.1 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncabatoff/go-seq v0.0.0-20180805175032-b08ef85ed833 // indirect
	github.com/nicolai86/scaleway-sdk v1.10.2-0.20180628010248-798f60e20bb2 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.96.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc5 // indirect
	github.com/opencontainers/runc v1.1.12 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.1 // indirect
	github.com/opencontainers/selinux v1.11.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.2 // indirect
	github.com/oschwald/maxminddb-golang v1.11.0
	github.com/ovh/go-ovh v1.4.3 // indirect
	github.com/packethost/packngo v0.1.1-0.20180711074735-b9cb5096f54c // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/prometheus-community/go-runit v0.1.0 // indirect
	github.com/prometheus/alertmanager v0.26.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0
	github.com/prometheus/exporter-toolkit v0.11.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/remeh/sizedwaitgroup v1.0.0 // indirect
	github.com/renier/xmlrpc v0.0.0-20170708154548-ce4a1a486c03 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/safchain/ethtool v0.3.0 // indirect
	github.com/samber/lo v1.38.1 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.21
	github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 // indirect
	github.com/seccomp/libseccomp-golang v0.9.2-0.20220502022130-f33da4d89646 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20230704072500-f1e31cf0ba5c // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 // indirect
	github.com/snowflakedb/gosnowflake v1.7.2-0.20240103203018-f1d625f17408 // indirect
	github.com/softlayer/softlayer-go v0.0.0-20180806151055-260589d94c7d // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.18.2
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.162 // indirect
	github.com/tg123/go-htpasswd v1.2.2 // indirect
	github.com/tilinna/clock v1.1.0
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/vertica/vertica-sql-go v1.3.3 // indirect
	github.com/vishvananda/netlink v1.2.1-beta.2 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/vmware/govmomi v0.36.1 // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/willf/bitset v1.1.11 // indirect
	github.com/willf/bloom v2.0.3+incompatible // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	github.com/xo/dburl v0.20.0 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.etcd.io/etcd/api/v3 v3.5.10 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.10 // indirect
	go.etcd.io/etcd/client/v3 v3.5.10 // indirect
	go.mongodb.org/mongo-driver v1.12.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.96.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.47.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel/bridge/opencensus v1.24.0 // indirect
	go4.org/netipx v0.0.0-20230125063823-8449b0a6169f // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	gonum.org/v1/gonum v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	howett.net/plist v1.0.0 // indirect
	k8s.io/kube-openapi v0.0.0-20231010175941-2dd684a91f00 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)

require github.com/ianlancetaylor/demangle v0.0.0-20240312041847-bd984b5ce465 // indirect

require (
	connectrpc.com/connect v1.14.0
	github.com/ethereum/go-ethereum v1.14.10
	github.com/ethpandaops/beacon v0.53.0
	github.com/ethpandaops/ethereum-metrics-exporter v0.27.4
	github.com/githubexporter/github-exporter v0.0.0-20231025122338-656e7dc33fe7
	github.com/grafana/agent-remote-config v0.0.2
	github.com/grafana/jfr-parser/pprof v0.0.3
	github.com/grafana/jsonparser v0.0.0-20240209175146-098958973a2d
	github.com/natefinch/atomic v1.0.1
	github.com/onrik/ethrpc v1.2.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.96.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.96.0
	go.opentelemetry.io/collector/config/configretry v0.96.0
	go.opentelemetry.io/collector/confmap/converter/expandconverter v0.96.0
	go.opentelemetry.io/collector/confmap/provider/fileprovider v0.96.0
	go.opentelemetry.io/collector/confmap/provider/yamlprovider v0.96.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.24.0
	golang.org/x/crypto/x509roots/fallback v0.0.0-20240208163226-62c9f1799c91
	k8s.io/apimachinery v0.29.2
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2 v2.2.1 // indirect
	github.com/DataDog/sketches-go v1.4.4 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.26.0 // indirect
	github.com/Shopify/sarama v1.38.1 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/Workiva/go-datastructures v1.1.0 // indirect
	github.com/attestantio/go-eth2-client v0.25.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/amp v1.23.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.22.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.19.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.38.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v1.36.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.147.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.20.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.27.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/shield v1.24.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/storagegateway v1.26.0 // indirect
	github.com/axiomhq/hyperloglog v0.0.0-20240124082744-24bca3a5b39b // indirect
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.4 // indirect
	github.com/channelmeter/iso8601duration v0.0.0-20150204201828-8da3af7a2a61 // indirect
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240223125850-b1e8a79f509c // indirect
	github.com/crate-crypto/go-kzg-4844 v1.0.0 // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/dgryski/go-metro v0.0.0-20180109044635-280f6062b5bc // indirect
	github.com/drone/envsubst v1.0.3 // indirect
	github.com/emicklei/dot v1.6.4 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.32.4 // indirect
	github.com/ethereum/c-kzg-4844 v1.0.0 // indirect
	github.com/ethereum/go-verkle v0.1.1-0.20240829091221-dffa7562dbe9 // indirect
	github.com/ethpandaops/ethwallclock v0.2.0 // indirect
	github.com/expr-lang/expr v1.16.1 // indirect
	github.com/ferranbt/fastssz v0.1.4 // indirect
	github.com/go-co-op/gocron v1.16.2 // indirect
	github.com/go-jose/go-jose/v4 v4.0.5 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/goccy/go-yaml v1.9.5 // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/grafana/jfr-parser v0.9.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hetznercloud/hcloud-go/v2 v2.4.0 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/huandu/go-clone v1.6.0 // indirect
	github.com/huin/goupnp v1.3.0 // indirect
	github.com/influxdata/tdigest v0.0.2-0.20210216194612-fc98d27c9e8b // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/knadh/koanf/v2 v2.1.0 // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/metalmatze/signal v0.0.0-20210307161603-1c9aa721a97a // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kafka v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/azure v0.96.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.96.0 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pk910/dynamic-ssz v0.0.4 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/prometheus-community/prom-label-proxy v0.6.0 // indirect
	github.com/prysmaticlabs/go-bitfield v0.0.0-20240618144021-706c95b2dd15 // indirect
	github.com/r3labs/sse/v2 v2.10.0 // indirect
	github.com/relvacode/iso8601 v1.4.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sercand/kuberesolver/v5 v5.1.1 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/skeema/knownhosts v1.2.1 // indirect
	github.com/sony/gobreaker v0.5.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/supranational/blst v0.3.13 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/envprovider v0.96.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/httpprovider v0.96.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/httpsprovider v0.96.0 // indirect
	go.opentelemetry.io/contrib/config v0.4.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.24.0 // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)

// NOTE: replace directives below must always be *temporary*.
//
// Adding a replace directive to change a module to a fork of a module will
// only be accepted when a PR upstream has been opened to accept the new
// change.
//
// Contributors are expected to work with upstream to make their changes
// acceptable, and remove the `replace` directive as soon as possible.
//
// If upstream is unresponsive, you should consider making a hard fork
// (i.e., creating a new Go module with the same source) or picking a different
// dependency.

// Replace directives from Prometheus
replace (
	k8s.io/klog => github.com/simonpasquier/klog-gokit v0.3.0
	// Prometheus uses v3.3.0, but we will get a compilation error from another module if we use it.
	k8s.io/klog/v2 => github.com/simonpasquier/klog-gokit/v3 v3.4.0
)

// TODO(marctc): remove replace directive once:
//
// * There is a release of Prometheus which contains
// prometheus/prometheus#13002
// and prometheus/prometheus#13497
// We use the last v1-related tag as the replace statement does not work for v2
// tags without the v2 suffix to the module root.
replace github.com/prometheus/prometheus => github.com/grafana/prometheus v1.8.2-0.20240130142130-51b39f24d406 // cmp_header_order branch

replace gopkg.in/yaml.v2 => github.com/rfratto/go-yaml v0.0.0-20211119180816-77389c3526dc

// Replace directives from Loki
replace (
	github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v36.2.0+incompatible
	github.com/Azure/azure-storage-blob-go => github.com/MasslessParticle/azure-storage-blob-go v0.14.1-0.20220216145902-b5e698eff68e
	github.com/bradfitz/gomemcache => github.com/themihai/gomemcache v0.0.0-20180902122335-24332e2d58ab
	github.com/go-kit/log => github.com/dannykopping/go-kit-log v0.2.2-0.20221002180827-5591c1641b6b
	github.com/gocql/gocql => github.com/grafana/gocql v0.0.0-20200605141915-ba5dc39ece85
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
	github.com/sercand/kuberesolver/v4 => github.com/sercand/kuberesolver/v5 v5.1.1
	github.com/thanos-io/thanos v0.22.0 => github.com/thanos-io/thanos v0.19.1-0.20211126105533-c5505f5eaa7d
	gopkg.in/Graylog2/go-gelf.v2 => github.com/grafana/go-gelf v0.0.0-20211112153804-126646b86de8
)

// TODO(rfratto): remove forks when changes are merged upstream
replace (
	// TODO(tpaschalis) this is to remove global instantiation of plugins
	// and allow non-singleton components.
	// https://github.com/grafana/cadvisor/tree/grafana-v0.47-noglobals
	github.com/google/cadvisor => github.com/grafana/cadvisor v0.0.0-20231110094609-5f7917925dea

	github.com/prometheus-community/postgres_exporter => github.com/grafana/postgres_exporter v0.8.1-0.20210722175051-db35d7c2f520

	// TODO(marctc): remove once this PR is merged upstream: https://github.com/prometheus/mysqld_exporter/pull/774
	github.com/prometheus/mysqld_exporter => github.com/grafana/mysqld_exporter v0.12.2-0.20231005125903-364b9c41e595

	// TODO(marctc, mattdurham): Replace node_export with custom fork for multi usage. https://github.com/prometheus/node_exporter/pull/2812
	github.com/prometheus/node_exporter => github.com/grafana/node_exporter v0.18.1-grafana-r01.0.20231004161416-702318429731
)

// Replacing for an internal fork which allows us to observe metrics produced by the Collector.
// This is a temporary solution while a new configuration design is discussed for the collector. Related issues:
// https://github.com/open-telemetry/opentelemetry-collector/issues/7532
// https://github.com/open-telemetry/opentelemetry-collector/pull/7644
// https://github.com/open-telemetry/opentelemetry-collector/pull/7696
// https://github.com/open-telemetry/opentelemetry-collector/issues/4970
replace (
	go.opentelemetry.io/collector/otelcol => github.com/grafana/opentelemetry-collector/otelcol v0.0.0-20240321103955-8919a1c85cbe
	go.opentelemetry.io/collector/service => github.com/grafana/opentelemetry-collector/service v0.0.0-20240321103955-8919a1c85cbe
)

// Required to avoid an ambiguous import with github.com/tencentcloud/tencentcloud-sdk-go
exclude github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.194

// Add exclude directives so Go doesn't pick old incompatible k8s.io/client-go
// versions.
exclude (
	k8s.io/client-go v8.0.0+incompatible
	k8s.io/client-go v12.0.0+incompatible
)

replace github.com/github/smimesign => github.com/grafana/smimesign v0.2.1-0.20220408144937-2a5adf3481d3

// Replacing for an internal for with a bugfix for delta histograms, https://github.com/grafana/stackdriver_exporter/pull/1
// Moving back to upstream is being tracked in an internal issue
replace github.com/prometheus-community/stackdriver_exporter => github.com/grafana/stackdriver_exporter v0.0.0-20240228143257-3a2c9acef5a2
