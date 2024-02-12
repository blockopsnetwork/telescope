package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type configs struct {
	config telescopeAttrs
}

type telescopeAttrs struct {
	ProjectName string `mapstructure:"projectName"`
	ProjectId   string `mapstructure:"projectId"`
	Network     string `mapstructure:"network"`
}

type Networks string

func (e *Networks) String() string {
	return string(*e)
}

// Check that the right parameters are passed
func (e *Networks) Set(v string) error {
	switch v {
	case "ethereum", "polkadot", "arbitrum", "base", "optimism":
		*e = Networks(v)
		return nil
	default:
		return errors.New(`invalid network. Please choose from: ethereum, polkadot, arbitrum, base, optimism`)
	}
}

func (e *Networks) Type() string {
	return "networks"
}

// Check and load telescope configuration to /home/.telescope.yaml
func LoadConfig() *configs {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".telescope")

	viper.AutomaticEnv()

	var c = &configs{}
	path := filepath.Join(home, ".telescope.yaml")

	if err := viper.ReadInConfig(); err == nil {
		log.Println("configuration file ", viper.ConfigFileUsed(), "found")
		if err := viper.Unmarshal(c); err != nil {
			cobra.CheckErr(err)
		}
	} else {
		log.Println("no configuration file found")

		file, err := os.Create(path)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer file.Close()
		fmt.Println("emtpy configuration file created at ", path)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("telescope")
	viper.BindEnv("config.projectName")
	viper.BindEnv("config.projectId")
	viper.BindEnv("config.network")
	viper.SafeWriteConfigAs(path)
	viper.WriteConfig()
	log.Println("successfully bind environment variables to configuration file values")
	return c
}

// Custom help command
func GetRootHelp() string {
	return `

Telescope is a one-stop observability tool for web3 observability.

Complete documentation is available at: https://docs.blockops.network/telescope

Synoposis:
	telescope [command] [flags]

Usage:
	telescope [flags]
	telescope [command]

Available Commands:
  help        Help about any command
  version     Print current telescope version
	monitor         Setting up telescope

Flags:
	Required:
	-p, --projectName          The name of the project
	-i, --projectId            The Id of the project
	-n, --network              Blockchain node network

	Optional:
	-c, --config               Config file (default is $HOME/.telescope
	-d, --debug                Outputs stack trace in case an exception is thrown
	-v, --version              Outputs release version


Examples:

	# run telescope
	telescope run -p calvin-trw -i 16-gdfsg -n polkadot

	# Enable debug mode
	telescope run -p calvin-trw -i 16-gdfsg -n polkadot -d

	# Print current version
	telescope version

Use "telescope [command] --help" for more information about a command.
`
}

// Set telescope configuration to a path
func SetTelescopeConfigs(projectName, projectId string, network Networks) bool {
	result := false
	if len(projectName) <= 0 {
		return result
	}

	viper.Set("config.projectName", projectName)
	viper.Set("config.projectId", projectId)
	viper.Set("config.network", network)
	viper.WriteConfig()
	result = true
	log.Println("telescope configs saved ", "to .telescope configuration file at ", viper.ConfigFileUsed())
	return result
}

// Verify the user suppplied all the required configs
func CheckConfigs(projectName, projectId string, network Networks) {
	savedProjectName := viper.GetString("config.projectName")
	savedProjectId := viper.GetString("config.projectId")
	savedNetwork := viper.GetString("config.network")
	if savedProjectName == "" || savedProjectId == "" || savedNetwork == "" {
		if projectName != "" && projectId != "" && network != "" {
			_ = SetTelescopeConfigs(projectName, projectId, network)
		} else if projectName == "" || projectId == "" || network == "" {
			fmt.Fprintln(os.Stderr, "--projectName|-p, --projectId|-i and --network|-n are required to run telescope")
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stderr, "--projectName|-p, --projectId|-i and --network|-n are required to run telescope")
			os.Exit(1)
		}
	}
}




