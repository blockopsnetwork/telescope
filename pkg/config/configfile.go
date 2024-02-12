package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/blockopsnetwork/telescope/internal/logger"
	"github.com/spf13/viper"
)

var (
	projectID = viper.GetString("config.projectId")
	projectName = viper.GetString("config.projectName")
	remote_write_url = "https://thanos-receiver.blockops.network/api/v1/receive"
	telescope_username = viper.GetString("config.telescope_username")
	telescope_password = viper.GetString("config.telescope_password")
)


func CreateConfigFile() {
	createConfigDir()
	logger.Log.Info("Creating configuration file...  ")
	telescopeConfigFile := fmt.Sprintf("%s/.telescope/agent.yaml", os.Getenv("HOME"))

	// if _, err := os.Stat(telescopeConfigFile); err == nil {
	// 	fmt.Println("Configuration file already exists. Skipping.")
	// 	return
	// }

	configContent := fmt.Sprintf(`server:
  log_level: info

metrics:
  wal_directory: /tmp/wal
  global:
    scrape_interval: 15s
    external_labels:
      project_id: %s
      project_name: %s
    remote_write:
      - url: %s
        basic_auth:
          username: %s
          password: %s

  configs:
`, projectID, projectName, remote_write_url, telescope_username, telescope_password)

  var Selectednetwork = viper.GetString("config.network")
	switch Selectednetwork {
	case "ethereum":
		configContent += `
    - name: geth
      scrape_configs:
        - job_name: geth
          static_configs:
            - targets: ["localhost:9100"]`
	case "polkadot":
		configContent += `
    - name: parachain
      scrape_configs:
        - job_name: parachain
          static_configs:
            - targets: ["localhost:9615"]
    - name: relaychain
      scrape_configs:
        - job_name: relaychain
          static_configs:
            - targets: ["localhost:9616"]`
	case "arbitrum":
		configContent += `
    - name: arbitrum
      scrape_configs:
        - job_name: arbitrum
          static_configs:
            - targets: ["localhost:9100"]`
	case "base":
		configContent += `
    - name: base
      scrape_configs:
        - job_name: base
          static_configs:
            - targets: ["localhost:9100"]`
	case "optimism":
		configContent += `
    - name: optimism
      scrape_configs:
        - job_name: optimism
          static_configs:
            - targets: ["localhost:9100"]`
	}

	configContent += `
integrations:
  agent:
    enabled: false
  node_exporter:
    enabled: true
`

	err := ioutil.WriteFile(telescopeConfigFile, []byte(configContent), 0644)
	if err != nil {
		logger.Log.Fatalf("Error creating configuration file:%v", err)
		return
	}

	logger.Log.Info("🎉 Configuration file created successfully.")
}

func createConfigDir() {
	configDir := fmt.Sprintf("%s/.telescope", os.Getenv("HOME"))
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			logger.Log.Fatalf("Error creating configuration directory:%v", err)
			return
		}
	}
}


