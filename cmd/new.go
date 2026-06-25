package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		if _, exists := s.Profiles[name]; exists {
			fatal(fmt.Sprintf("profile %q already exists", name))
		}
		s.Profiles[name] = store.Profile{Vars: make(map[string]string)}
		if err := s.Save(); err != nil {
			fatal(err.Error())
		}
		printOK(fmt.Sprintf("Created profile %q", name))
		if s.Active == "" {
			printInfo(fmt.Sprintf("  Run: w3env use %s", name))
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
