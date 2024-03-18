#!/bin/bash

set -eo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
BGREEN='\033[1;32m'
AMBER='\033[0;33m'
BWHITE='\033[1;37m'
NC='\033[0m'

TELESCOPE_IMAGE=grafana/agent:v0.37.2
TELESCOPE_DIR="${HOME}/.telescope"


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
        --telescope-username)
            TELESCOPE_USERNAME=$2
            shift 2;;
        --telescope-password)
            TELESCOPE_PASSWORD=$2
            shift 2;;
        --remote-write-url)
            REMOTE_WRITE_URL=$2
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
    mkdir -p "${TELESCOPE_DIR}"
    if [[ 0 -ne $? ]]
    then
        error "Failed to create the configuration directory."
    fi

    info "OK"
}


# Create agent configuration file
function create_config_file () {
    create_config_dir
    echo -n "Creating configuration file...  "
    TELESCOPE_CONFIG_FILE="${HOME}/.telescope/agent.yaml"
    
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
      - url: ${REMOTE_WRITE_URL}
        basic_auth:
          username: ${TELESCOPE_USERNAME}
          password: ${TELESCOPE_PASSWORD}

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
    docker pull --quiet "${TELESCOPE_IMAGE}"
    if [[ 0 -ne $? ]]
    then
        error "Failed to pull the telescope docker image."
    fi

    info "OK"
}

function run_telescope () {
    echo -n "Setting up Telescope...  "
    create_config_file
    get_docker_image
    docker run -d  \
        --name telescope --network="host" --pid="host" --cap-add SYS_TIME \
        -v "${TELESCOPE_DIR}":/etc/agent-config \
        --restart unless-stopped \
        -e PROJECT_ID=${PROJECT_ID} \
        -e PROJECT_NAME=${PROJECT_NAME} \
        -e NETWORK=${NETWORK} \
        -e REMOTE_WRITE_URL=${REMOTE_WRITE_URL} \
        -e TELESCOPE_USERNAME=${TELESCOPE_USERNAME} \
        -e TELESCOPE_PASSWORD=${TELESCOPE_PASSWORD} \
        "${TELESCOPE_IMAGE}" --config.file=/etc/agent-config/agent.yaml -config.expand-env  
    if [[ 0 -ne $? ]]
    then
        error "Failed to create the docker container."
    fi

    success_bold "Telescope setup completed successfully for ${NETWORK} network."
    info "$(docker logs -f telescope)"

}

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
