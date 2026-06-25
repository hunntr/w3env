package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var setCmd = &cobra.Command{
	Use:   "set <key> <value>  |  set KEY=value",
	Short: "Set a variable in the active profile",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var key, value string
		if len(args) == 1 {
			parts := strings.SplitN(args[0], "=", 2)
			if len(parts) != 2 {
				fatal("usage: w3env set <key> <value>  or  w3env set KEY=value")
			}
			key, value = parts[0], parts[1]
		} else {
			key, value = args[0], args[1]
		}
		if key == "" {
			fatal("key cannot be empty")
		}
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		if s.Active == "" {
			fatal("no active profile - run: w3env use <profile>")
		}
		p := s.Profiles[s.Active]
		if p.Vars == nil {
			p.Vars = make(map[string]string)
		}
		p.Vars[key] = value
		s.Profiles[s.Active] = p
		if err := s.Save(); err != nil {
			fatal(err.Error())
		}
		writeActivation(s, s.Active)
		printOK(fmt.Sprintf("Set %s in %q", col(cBold, key), s.Active))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
