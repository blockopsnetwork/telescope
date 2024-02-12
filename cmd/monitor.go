/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/blockopsnetwork/telescope/internal/process"
	"github.com/blockopsnetwork/telescope/pkg/config"
	"github.com/blockopsnetwork/telescope/pkg/text"
	"github.com/spf13/cobra"
)

var (
	projectName string
	projectId   string
	network     config.Networks
)

var runCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Setting up telescope",
	Long:  text.AdditionalText,
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfigs(projectName, projectId, network)
		config.CreateConfigFile()
		cl := process.NewDockerClient()
		cl.PullImage()
		cl.StartDockerContainer()
		cl.ContainerLogs()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&projectName, "projectName", "p", "", "The name of the project")
	runCmd.Flags().StringVarP(&projectId, "projectId", "i", "", "The Id of the project")
	runCmd.Flags().VarP(&network, "network", "n", `Network specific configuration. allowed: "arbitrum", "ethereum", "polkadot", "base", "optimism"`)
	runCmd.RegisterFlagCompletionFunc("network", networkCompletion)
}

func networkCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"arbitrum\thelp text for arbitrum",
		"ethereum\thelp text for ethereum",
		"polkadot\thelp text for polkadot",
		"base\thelp text for base",
		"optimism\thelp text for optimism",
	}, cobra.ShellCompDirectiveDefault
}
