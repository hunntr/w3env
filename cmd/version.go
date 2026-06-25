package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.1.0-beta"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("w3env " + Version)
	},
}

func init() {
	rootCmd.Version = Version
	rootCmd.AddCommand(versionCmd)
}
