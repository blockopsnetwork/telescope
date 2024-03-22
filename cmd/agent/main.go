package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-kit/log/level"
	"github.com/grafana/agent/internal/build"



func init() {
	prometheus.MustRegister(build.NewCollector("agent"))
}

func main() {
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