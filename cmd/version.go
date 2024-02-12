package cmd

import (
	"fmt"

	"github.com/blockopsnetwork/telescope/pkg/text"
	"github.com/spf13/cobra"
)

var (
	build       = "0"
	commit      = "sha"
	releaseDate = "2024-02-01"
	version     = "v0.1.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current telescope version",
	Long:  text.AdditionalText,
	Run: func(cmd *cobra.Command, args []string) {
		sh, err := cmd.Flags().GetBool("short")
		cobra.CheckErr(err)
		if sh {
			fmt.Println("Telescope CLI " + version)
		} else {
			fmt.Println("Telescope CLI", version)
			fmt.Println("Build", build)
			fmt.Println("Release Date", releaseDate)
			fmt.Println("Commit", commit)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolP("short", "s", false, "short discription of telescope API version")
}
