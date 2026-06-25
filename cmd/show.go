package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var showReveal bool

var showCmd = &cobra.Command{
	Use:     "show [profile]",
	Aliases: []string{"env"},
	Short:   "Show all variables in a profile (active profile if none specified)",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}

		var name string
		var p *store.Profile

		if len(args) == 1 {
			name = args[0]
			prof, exists := s.Profiles[name]
			if !exists {
				fatal(fmt.Sprintf("profile %q not found", name))
			}
			p = &prof
		} else {
			prof, err := s.ActiveProfile()
			if err != nil {
				fatal(err.Error())
			}
			p = prof
			name = s.Active
		}

		label := name
		if name == s.Active {
			label = name + col(cGreen, " (active)")
		}
		fmt.Printf("Profile: %s\n\n", col(cBold+cCyan, label))

		if len(p.Vars) == 0 {
			printInfo("  (empty - add vars with: w3env set <key> <value>)")
			return
		}

		keys := make([]string, 0, len(p.Vars))
		maxLen := 0
		for k := range p.Vars {
			keys = append(keys, k)
			if len(k) > maxLen {
				maxLen = len(k)
			}
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := p.Vars[k]
			pad := strings.Repeat(" ", maxLen-len(k))
			display := v
			if !showReveal && isSensitive(k) {
				display = col(cGray, "****** (use --reveal to show)")
			}
			fmt.Printf("  %s%s  =  %s\n", col(cCyan, k), pad, display)
		}
		fmt.Printf("\n%s\n", col(cGray, fmt.Sprintf("%d variable(s)", len(p.Vars))))
	},
}

func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, pat := range []string{"key", "secret", "password", "pass", "token", "pk", "mnemonic", "seed"} {
		if strings.Contains(lower, pat) {
			return true
		}
	}
	return false
}

func init() {
	showCmd.Flags().BoolVar(&showReveal, "reveal", false, "Show sensitive values in plain text")
	rootCmd.AddCommand(showCmd)
}
