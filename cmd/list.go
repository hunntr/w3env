package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all profiles",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		names := s.ProfileNames()
		if len(names) == 0 {
			printInfo("No profiles yet - run: w3env new <name>")
			return
		}
		for _, n := range names {
			p := s.Profiles[n]
			if n == s.Active {
				fmt.Printf("  %s  %-24s %s\n",
					col(cGreen, "▶"),
					col(cBold, n),
					col(cGray, fmt.Sprintf("(%d vars)", len(p.Vars))),
				)
			} else {
				fmt.Printf("     %-24s %s\n",
					n,
					col(cGray, fmt.Sprintf("(%d vars)", len(p.Vars))),
				)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
