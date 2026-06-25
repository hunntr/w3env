package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate the current profile and unset its variables",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		if s.Active == "" {
			printWarn("no active profile")
			return
		}
		name := s.Active
		writeDeactivation(s, name)
		s.Active = ""
		if err := s.Save(); err != nil {
			fatal(err.Error())
		}
		fmt.Fprintf(os.Stderr, "%s Deactivated %q\n", col(cGray, "○"), name)
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)
}
