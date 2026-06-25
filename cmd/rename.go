package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a profile",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName, newName := args[0], args[1]
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		p, exists := s.Profiles[oldName]
		if !exists {
			fatal(fmt.Sprintf("profile %q not found", oldName))
		}
		if oldName == newName {
			return
		}
		if _, exists := s.Profiles[newName]; exists {
			fatal(fmt.Sprintf("profile %q already exists", newName))
		}
		s.Profiles[newName] = p
		delete(s.Profiles, oldName)
		if s.Active == oldName {
			s.Active = newName
		}
		if err := s.Save(); err != nil {
			fatal(err.Error())
		}
		printOK(fmt.Sprintf("Renamed %q -> %q", oldName, newName))
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
