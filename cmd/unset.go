package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var unsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Remove a variable from the active profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		if s.Active == "" {
			fatal("no active profile - run: w3env use <profile>")
		}
		p := s.Profiles[s.Active]
		if _, exists := p.Vars[key]; !exists {
			fatal(fmt.Sprintf("key %q not found in profile %q", key, s.Active))
		}
		delete(p.Vars, key)
		s.Profiles[s.Active] = p
		if err := s.Save(); err != nil {
			fatal(err.Error())
		}
		writeActivation(s, s.Active)
		printOK(fmt.Sprintf("Removed %s from %q", col(cBold, key), s.Active))
	},
}

func init() {
	rootCmd.AddCommand(unsetCmd)
}
