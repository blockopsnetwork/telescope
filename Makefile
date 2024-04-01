# include tools/make/*.mk

AGENT_IMAGE                             ?= blockopsnetwork/agent:latest
OPERATOR_IMAGE                          ?= blockopsnetwork/agent:latest
AGENT_BINARY                            ?= build/agent
OPERATOR_BINARY                         ?= build/agent-operator
AGENTLINT_BINARY                        ?= build/agentlint
GOOS                                    ?= $(shell go env GOOS)
GOARCH                                  ?= $(shell go env GOARCH)
GOARM                                   ?= $(shell go env GOARM)
CGO_ENABLED                             ?= 1
RELEASE_BUILD                           ?= 0
GOEXPERIMENT                            ?= $(shell go env GOEXPERIMENT)

# List of all environment variables which will propagate to the build
# container. USE_CONTAINER must _not_ be included to avoid infinite recursion.
PROPAGATE_VARS := \
    AGENT_IMAGE OPERATOR_IMAGE \
    BUILD_IMAGE GOOS GOARCH GOARM CGO_ENABLED RELEASE_BUILD \
    AGENT_BINARY OPERATOR_BINARY \
    VERSION GO_TAGS GOEXPERIMENT

#
# Constants for targets
#

GO_ENV := GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) CGO_ENABLED=$(CGO_ENABLED)

VERSION      ?= $(shell bash ./tools/image-tag)
GIT_REVISION := $(shell git rev-parse --short HEAD)
GIT_BRANCH   := $(shell git rev-parse --abbrev-ref HEAD)
VPREFIX      := github.com/blockopsnetwork/telescope/internal/build
GO_LDFLAGS   := -X $(VPREFIX).Branch=$(GIT_BRANCH)                        \
                -X $(VPREFIX).Version=$(VERSION)                          \
                -X $(VPREFIX).Revision=$(GIT_REVISION)                    \
                -X $(VPREFIX).BuildUser=$(shell whoami)@$(shell hostname) \
                -X $(VPREFIX).BuildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

DEFAULT_FLAGS    := $(GO_FLAGS)
DEBUG_GO_FLAGS   := -ldflags "$(GO_LDFLAGS)" -tags "$(GO_TAGS)"
RELEASE_GO_FLAGS := -ldflags "-s -w $(GO_LDFLAGS)" -tags "$(GO_TAGS)"

ifeq ($(RELEASE_BUILD),1)
GO_FLAGS := $(DEFAULT_FLAGS) $(RELEASE_GO_FLAGS)
else
GO_FLAGS := $(DEFAULT_FLAGS) $(DEBUG_GO_FLAGS)
endif

#
# Targets for running tests
#
# These targets currently don't support proxying to a build container due to
# difficulties with testing ./internal/util/k8s and testing packages.
#

.PHONY: lint
lint: agentlint
	golangci-lint run -v --timeout=10m
	$(AGENTLINT_BINARY) ./...

.PHONY: test
# We have to run test twice: once for all packages with -race and then once
# more without -race for packages that have known race detection issues.
test:
	$(GO_ENV) go test $(GO_FLAGS) -race $(shell go list ./... | grep -v /integration-tests/)
	$(GO_ENV) go test $(GO_FLAGS) ./internal/static/integrations/node_exporter ./internal/static/logs ./internal/static/operator ./internal/util/k8s ./internal/component/otelcol/processor/tail_sampling ./internal/component/loki/source/file ./internal/component/loki/source/docker

test-packages:
	docker pull $(BUILD_IMAGE)
	go test -tags=packaging  ./internal/tools/packaging_test

.PHONY: integration-test
integration-test:
	cd internal/cmd/integration-tests && $(GO_ENV) go run .

#
# Targets for building binaries
#

.PHONY: binaries agent operator
binaries: agent operator

agent:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	$(GO_ENV) go build $(GO_FLAGS) -o $(AGENT_BINARY) ./cmd/agent
endif

operator:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	$(GO_ENV) go build $(GO_FLAGS) -o $(OPERATOR_BINARY) ./cmd/agent-operator
endif

agentlint:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	cd ./internal/cmd/agentlint && $(GO_ENV) go build $(GO_FLAGS) -o ../../../$(AGENTLINT_BINARY) .
endif

#
# Targets for building Docker images
#

DOCKER_FLAGS := --build-arg RELEASE_BUILD=$(RELEASE_BUILD) --build-arg VERSION=$(VERSION)

ifneq ($(DOCKER_PLATFORM),)
DOCKER_FLAGS += --platform=$(DOCKER_PLATFORM)
endif

.PHONY: images agent-image operator-image
images: agent-image operator-image

agent-image:
	DOCKER_BUILDKIT=1 docker build $(DOCKER_FLAGS) -t $(AGENT_IMAGE) -f cmd/agent/Dockerfile .
operator-image:
	DOCKER_BUILDKIT=1 docker build $(DOCKER_FLAGS) -t $(OPERATOR_IMAGE) -f cmd/agent-operator/Dockerfile .

#
# Targets for generating assets
#

.PHONY: generate generate-crds generate-drone generate-helm-docs generate-helm-tests generate-dashboards generate-protos generate-ui generate-versioned-files
generate: generate-crds generate-drone generate-helm-docs generate-helm-tests generate-dashboards generate-protos generate-ui generate-versioned-files generate-docs

generate-crds:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	bash ./tools/generate-crds.bash
	gen-crd-api-reference-docs -config tools/gen-crd-docs/config.json -api-dir "github.com/blockopsnetwork/telescope/internal/static/operator/apis/monitoring/" -out-file docs/sources/operator/api.md -template-dir tools/gen-crd-docs/template
endif

generate-helm-docs:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	cd operations/helm/charts/grafana-agent && helm-docs
endif

generate-helm-tests:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	bash ./operations/helm/scripts/rebuild-tests.sh
endif

generate-dashboards:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	cd example/docker-compose && jb install && \
	cd grafana/dashboards && jsonnet template.jsonnet -J ../../vendor -m .
endif

generate-protos:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	go generate ./internal/static/agentproto/
endif

generate-ui:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	cd ./internal/web/ui && yarn --network-timeout=1200000 && yarn run build
endif

generate-versioned-files:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	sh ./tools/gen-versioned-files/gen-versioned-files.sh
endif

generate-docs:
ifeq ($(USE_CONTAINER),1)
	$(RERUN_IN_CONTAINER)
else
	go generate ./docs
endif