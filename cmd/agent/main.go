package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-kit/log/level"
	"github.com/grafana/agent/internal/boringcrypto"
	"github.com/grafana/agent/internal/build"
	"github.com/grafana/agent/internal/flowmode"
	"github.com/grafana/agent/internal/static/config"
	"github.com/grafana/agent/internal/static/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// tconfig "github.com/grafana/agent/internal/telescope"
	util_log "github.com/grafana/agent/internal/util/log"

	"github.com/prometheus/client_golang/prometheus"

	// Register Prometheus SD components
	_ "github.com/prometheus/prometheus/discovery/consul"
	_ "github.com/prometheus/prometheus/discovery/install"

	// Register integrations
	_ "github.com/grafana/agent/internal/static/integrations/install"

	// Embed a set of fallback X.509 trusted roots
	// Allows the app to work correctly even when the OS does not provide a verifier or systems roots pool
	_ "golang.org/x/crypto/x509roots/fallback"
)

var cfgFile string

type TelescopeConfig struct {
	Metrics           bool
	Network           string
	ProjectId         string
	ProjectName       string
	TelescopeUsername string
	TelescopePassword string
	RemoteWriteUrl    string
}

var cmd = &cobra.Command{
	Use:   "metrics",
	Short: "A Web3 Observability tooling",
	Long:  `Gain insights into your Web3 infrastructure with Telescope.`,
	Run: func(cmd *cobra.Command, args []string) {
		var config TelescopeConfig
		if err := config.loadConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Telescope Config: %v\n", config)
	},
}

var helperFunction = func(cmd *cobra.Command, args []string) {
	fmt.Println("Custom Help:")
	fmt.Println("Usage: app start [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  --metrics\t\tEnable metrics")
	fmt.Println("  --network\t\tSpecify the network")
	fmt.Println("  --project-id\t\tSpecify the project ID")
	fmt.Println("  --project-name\tSpecify the project name")
	fmt.Println("  --telescope-username\tSpecify the telescope username")
	fmt.Println("  --telescope-password\tSpecify the telescope password")
	fmt.Println("  --remote-write-url\tSpecify the remote write URL")
}

func init() {
	prometheus.MustRegister(build.NewCollector("agent"))
	cobra.OnInitialize(initConfig)
	cmd.SetHelpFunc(helperFunction)
	cmd.Flags().Bool("metrics", false, "Enable metrics")
	cmd.Flags().String("network", "", "Specify the network")
	cmd.Flags().String("project-id", "", "Specify the project ID")
	cmd.Flags().String("project-name", "", "Specify the project name")
	cmd.Flags().String("telescope-username", "", "Specify the telescope username")
	cmd.Flags().String("telescope-password", "", "Specify the telescope password")
	cmd.Flags().String("remote-write-url", "", "Specify the remote write URL")

	// Bind flags with viper
	viper.BindPFlag("metrics", cmd.Flags().Lookup("metrics"))
	viper.BindPFlag("network", cmd.Flags().Lookup("network"))
	viper.BindPFlag("project-id", cmd.Flags().Lookup("project-id"))
	viper.BindPFlag("project-name", cmd.Flags().Lookup("project-name"))
	viper.BindPFlag("telescope-username", cmd.Flags().Lookup("telescope-username"))
	viper.BindPFlag("telescope-password", cmd.Flags().Lookup("telescope-password"))
	viper.BindPFlag("remote-write-url", cmd.Flags().Lookup("remote-write-url"))

	tconfig.generateConfigFile(config * TelescopeConfig)

	// cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cmd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func (c *TelescopeConfig) loadConfig() error {
	c.Metrics = viper.GetBool("metrics")
	c.Network = viper.GetString("network")
	c.ProjectId = viper.GetString("project-id")
	c.ProjectName = viper.GetString("project-name")
	c.TelescopeUsername = viper.GetString("telescope-username")
	c.TelescopePassword = viper.GetString("telescope-password")
	c.RemoteWriteUrl = viper.GetString("remote-write-url")

	// Check for required fields
	if err := checkRequiredFlags(); err != nil {
		return err
	}

	return nil

}

// func init() {
// 	prometheus.MustRegister(build.NewCollector("agent"))
// }

func agent() {
	runMode, err := getRunMode()
	if err != nil {
		log.Fatalln(err)
	}

	if runMode == runModeFlow {
		flowmode.Run()
		return
	}

	// Set up logging using default values before loading the config
	defaultCfg := server.DefaultConfig()
	logger := server.NewLogger(&defaultCfg)

	reloader := func(log *server.Logger) (*config.Config, error) {
		fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		return config.Load(fs, os.Args[1:], log)
	}
	cfg, err := reloader(logger)
	if err != nil {
		log.Fatalln(err)
	}

	// After this point we can start using go-kit logging.
	logger = server.NewLogger(cfg.Server)
	util_log.Logger = logger

	level.Info(logger).Log("boringcrypto enabled", boringcrypto.Enabled)
	ep, err := NewEntrypoint(logger, cfg, reloader)
	if err != nil {
		level.Error(logger).Log("msg", "error creating the agent server entrypoint", "err", err)
		os.Exit(1)
	}

	if err = ep.Start(); err != nil {
		level.Error(logger).Log("msg", "error running agent", "err", err)
		// Don't os.Exit here; we want to do cleanup by stopping promMetrics
	}

	ep.Stop()
	level.Info(logger).Log("msg", "agent exiting")
}

func main() {
	err := cmd.Execute()
	agent()
	if err != nil {
		os.Exit(1)
	}
}

func checkRequiredFlags() error {
	requiredFlags := []string{"metrics", "network", "project-id", "project-name", "telescope-username", "telescope-password", "remote-write-url"}
	missingFlags := []string{}

	for _, flag := range requiredFlags {
		if viper.GetString(flag) == "" {
			missingFlags = append(missingFlags, flag)
		}
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("Error: Missing required flags: %s", strings.Join(missingFlags, ", "))
	}

	return nil
}

func createConfigFIle(config *TelescopeConfig) error {
	telescopeDir := "/tmp/telescope"
	telescopeConfigFile := filepath.Join(telescopeDir, "telescope.yaml")

	// Create the configuration directory if it doesn't exist
	if _, err := os.Stat(telescopeDir); os.IsNotExist(err) {
		if err := os.Mkdir(telescopeDir, 0755); err != nil {
			return fmt.Errorf("Error creating configuration directory: %v", err)
		}
	}

	fmt.Print("Creating configuration file...  ")

	// Check if the configuration file already exists
	if _, err := os.Stat(telescopeConfigFile); err == nil {
		fmt.Println("Configuration file already exists. Skipping.")
		return nil
	}

	// Create the configuration file
	file, err := os.Create(telescopeConfigFile)
	if err != nil {
		return fmt.Errorf("Error creating configuration file: %v", err)
	}
	defer file.Close()

}
