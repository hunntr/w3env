package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Print the value of a variable (raw, for scripting)",
	Long: `Print the raw value of a variable - useful in scripts:
  cast call $(w3env get TARGET) "owner()(address)"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		p, err := s.ActiveProfile()
		if err != nil {
			fatal(err.Error())
		}
		v, ok := p.Vars[key]
		if !ok {
			fmt.Fprintf(os.Stderr, col(cRed, "✗")+" key %q not found in profile %q\n", key, s.Active)
			os.Exit(1)
		}
		fmt.Println(v)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
