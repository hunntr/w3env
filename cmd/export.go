package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Print shell export commands for the active profile",
	Long: `Print export commands for the active profile. Use in scripts:
  eval $(w3env export)
  source <(w3env export)`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		p, err := s.ActiveProfile()
		if err != nil {
			fatal(err.Error())
		}
		for k, v := range p.Vars {
			fmt.Printf("export %s=%q\n", k, v)
		}
		fmt.Printf("export W3ENV_ACTIVE=%q\n", s.Active)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
