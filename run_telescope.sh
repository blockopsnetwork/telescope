#!/bin/bash

set -eo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
BGREEN='\033[1;32m'
AMBER='\033[0;33m'
BWHITE='\033[1;37m'
NC='\033[0m'


PROJECT_ID=$1
PROJECT_NAME=$2
NETWORK=$3


function error() {
    echo -e "\n${RED}$*${NC}"
    exit 1
}

function success() {
    echo -e "\n${GREEN}$*${NC}"
}

function success_bold() {
    echo -e "\n${BGREEN}$*${NC}"
}

function warning() {
    echo -e "\n${AMBER}$*${NC}"
}

function info() {
    echo -e "\n${BWHITE}$*${NC}"
}

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            help
            exit 0
            ;;
        --metrics) 
            METRICS=true
            shift;;
        --network)
            NETWORK=$2
            shift 2;;
        --project-id)
            PROJECT_ID=$2
            shift 2;;
        --project-name)
            PROJECT_NAME=$2
            shift 2;;
        *)
            break
            ;;
    esac
done

# Network specific configuration
case "${NETWORK}" in
    ethereum)
        NETWORK="ethereum"
        SCRAPE_PORT=9100
        ;;
    polkadot)
        NETWORK="polkadot"
        RELAYCHAIN_PORT=9616
        PARACHAIN_PORT=9615
        ;;
    arbitrum)
        NETWORK="arbitrum"
        ;;
    base)
        NETWORK="base"
        ;;
    optimism)
        NETWORK="optimism"
        ;;
    *)
        error "Invalid network. Please choose from: ethereum, polkadot, arbitrum, base, optimism"
        ;;
esac

# Check that the right parameters are passed
if [ -z "${PROJECT_ID}" ] || [ -z "${PROJECT_NAME}" ] || [ -z "${NETWORK}" ]; then
    error "Usage: $0 <PROJECT_ID> <PROJECT_NAME> <NETWORK>\nNETWORK options: ethereum, polkadot, arbitrum, base, optimism"
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    error "Docker could not be found. Please install Docker and retry."
fi

function create_config_dir () {
    echo -n "Creating configuration directory...  "
    mkdir -p "${HOME}/.telescope"
    if [[ 0 -ne $? ]]
    then
        error "Failed to create the configuration directory."
    fi

    info "OK"
}


# Create agent configuration file
function create_config_file () {
    echo -n "Creating configuration file...  "
    TELESCOPE_CONFIG_FILE="${HOME}/.telescope/config.yaml"
    
    if [[ -f "${TELESCOPE_CONFIG_FILE}" ]]; then
        warning "Configuration file already exists. Skipping."
    else
        cat <<EOF > "${TELESCOPE_CONFIG_FILE}"
server:
  log_level: info

metrics:
  wal_directory: /tmp/wal
  global:
    scrape_interval: 15s
    external_labels:
      project_id: ${PROJECT_ID}
      project_name: ${PROJECT_NAME}
    remote_write:
      - url: https://prometheus-us-central1.grafana.net/api/prom/push
        basic_auth:
          username: 653537
          password: glc_eyJvIjoiNzQ1MTAwIiwibiI6InN0YWNrLTQ3NjM4Mi1obS13cml0ZS10ZXRzIiwiayI6ImFwaDFPMThqTjI4T0wzSzRoME04SG1mOSIsIm0iOnsiciI6InVzIn19

  configs:
EOF
        case "${NETWORK}" in
            ethereum)
                cat <<EOF >> "${TELESCOPE_CONFIG_FILE}"
    - name: geth
      scrape_configs:
        - job_name: geth
          static_configs:
            - targets: ["localhost:9100"]
EOF
                ;;
            polkadot)
                cat <<EOF >> "${TELESCOPE_CONFIG_FILE}"
    - name: parachain
      scrape_configs:
        - job_name: parachain
          static_configs:
            - targets: ["localhost:9615"]
    - name: relaychain
      scrape_configs:
        - job_name: relaychain
          static_configs:
            - targets: ["localhost:9616"]
EOF
                ;;
            arbitrum)
                cat <<EOF >> "${TELESCOPE_CONFIG_FILE}"
    - name: arbitrum
      scrape_configs:
        - job_name: arbitrum
          static_configs:
            - targets: ["localhost:9100"]
EOF
                ;;
            base)
                cat <<EOF >> "${TELESCOPE_CONFIG_FILE}"
    - name: base
        scrape_configs:
            - job_name: base
                static_configs:
                - targets: ["localhost:9100"]
EOF
                ;;
            optimism)
                cat <<EOF >> "${TELESCOPE_CONFIG_FILE}"
    - name: optimism
        scrape_configs:
            - job_name: optimism
                static_configs:
                - targets: ["localhost:9100"]
EOF
                ;;
        esac
        cat <<EOF >> "${TELESCOPE_CONFIG_FILE}"
integrations:
  agent:
    enabled: false
  node_exporter:
    enabled: true
EOF
        success "Configuration file created successfully."
    fi
}

function get_docker_image () {
    echo -n "Pulling docker image...  "
    TELESCOPE_IMAGE=grafana/agent:v0.37.2
    docker pull --quiet "${TELESCOPE_IMAGE}"
    if [[ 0 -ne $? ]]
    then
        error "Failed to pull the telescope docker image."
    fi

    info "OK"
}

function create_docker_volume () {
    echo -n "Creating docker volume...  "
    TELESCOPE_VOLUME=telescope
    docker volume create "${TELESCOPE_VOLUME}" &> /dev/null
    if [[ 0 -ne $? ]]
    then
        warning "Docker volume already exists. Skipping."
    fi

    info "OK"
}

function run_telescope () {
    echo -n "Setting up Telescope...  "
    TELESCOPE_CONTAINER=telescope
    docker run -d --rm \
        --name "host" --network="host" --pid="host" --cap-add SYS_TIME \
        --mount "type=volume,source=${TELESCOPE_VOLUME},destination=/var/lib/grafana-agent" \
        --restart unless-stopped \
        -e PROJECT_ID=${PROJECT_ID} \
        -e PROJECT_NAME=${PROJECT_NAME} \
        "${TELESCOPE_IMAGE}" --config.file=/etc/agent-config/agent.yaml -config.expand-env  &> /dev/null
    if [[ 0 -ne $? ]]
    then
        error "Failed to create the docker container."
    fi

    success_bold "Telescope setup completed successfully for ${NETWORK} network."

}

create_config_dir
create_config_file
get_docker_image
create_docker_volume
run_telescope


function help() {
    echo -e "${BWHITE}Usage:${NC} $0 <PROJECT_ID> <PROJECT_NAME> <NETWORK>"
    echo -e "${BWHITE}Where:${NC}"
    echo -e "  ${GREEN}PROJECT_ID${NC}    is your project identifier."
    echo -e "  ${GREEN}PROJECT_NAME${NC}  is the name of your project."
    echo -e "  ${GREEN}NETWORK${NC}       is the blockchain network for which Telescope is being set up."
    echo -e "${BWHITE}Network options:${NC}"
    echo -e "  ${AMBER}ethereum${NC}   - Set up Telescope on Ethereum Node."
    echo -e "  ${AMBER}polkadot${NC}   - Set up Telescope on Polkadot Node."
    echo -e "  ${AMBER}arbitrum${NC}   - Set up Telescope on Arbitrum Node."
    echo -e "  ${AMBER}base${NC}       - Set up Telescope on Base Node."
    echo -e "  ${AMBER}optimism${NC}   - Set up Telescope on Optimism Node."
    echo -e "\n${BWHITE}Example:${NC}"
    echo -e "  $0 my_project_id my_project_name ethereum"
}
