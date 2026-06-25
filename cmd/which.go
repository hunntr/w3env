package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var whichCmd = &cobra.Command{
	Use:   "which",
	Short: "Print the name of the active profile",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		if s.Active == "" {
			fmt.Fprintln(os.Stderr, "no active profile")
			os.Exit(1)
		}
		fmt.Println(s.Active)
	},
}

func init() {
	rootCmd.AddCommand(whichCmd)
}
