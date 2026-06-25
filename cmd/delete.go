package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:     "delete <name>",
	Aliases: []string{"rm", "del"},
	Short:   "Delete a profile",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		if _, exists := s.Profiles[name]; !exists {
			fatal(fmt.Sprintf("profile %q not found", name))
		}
		if !deleteForce {
			fmt.Printf("Delete profile %q? [y/N] ", col(cBold, name))
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
				fmt.Println("Aborted.")
				return
			}
		}
		if s.Active == name {
			writeDeactivation(s, name)
			s.Active = ""
		}
		delete(s.Profiles, name)
		if err := s.Save(); err != nil {
			fatal(err.Error())
		}
		printOK(fmt.Sprintf("Deleted profile %q", name))
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(deleteCmd)
}
