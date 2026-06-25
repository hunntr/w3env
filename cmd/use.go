package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var useCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Activate a profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		p, exists := s.Profiles[name]
		if !exists {
			fatal(fmt.Sprintf("profile %q not found - run: w3env new %s", name, name))
		}
		s.Active = name
		if err := s.Save(); err != nil {
			fatal(err.Error())
		}
		writeActivation(s, name)
		plural := "s"
		if len(p.Vars) == 1 {
			plural = ""
		}
		fmt.Fprintf(os.Stderr, "%s Switched to %s%s%s (%d var%s)\n",
			col(cGreen, "✓"), cBold, name, cReset, len(p.Vars), plural)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
