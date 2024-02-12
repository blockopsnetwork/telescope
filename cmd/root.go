/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/blockopsnetwork/telescope/pkg/config"
	"github.com/blockopsnetwork/telescope/pkg/text"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "telescope",
	Short: "Node Monitoring Platform for blockchain nodes",
	Long:  text.RootText,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stdout, text.AdditionalText)
		fmt.Fprintln(os.Stdout, "Use telescope --help for more information on how to use this CLI")
		os.Exit(0)
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableSuggestions = false
	rootCmd.PersistentFlags().SortFlags = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "c", "config file (default is $HOME/.telescope.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enables debug mode")

	// rootCmd.SetHelpTemplate(config.GetRootHelp())
}

func initConfig() {
	if debug {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(io.Discard)
	}
	config.LoadConfig()
}
