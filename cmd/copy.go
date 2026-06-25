package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var copyCmd = &cobra.Command{
	Use:     "copy <src> <dst>",
	Aliases: []string{"cp"},
	Short:   "Copy a profile",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src, dst := args[0], args[1]
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		srcProfile, exists := s.Profiles[src]
		if !exists {
			fatal(fmt.Sprintf("profile %q not found", src))
		}
		if _, exists := s.Profiles[dst]; exists {
			fatal(fmt.Sprintf("profile %q already exists", dst))
		}
		newProfile := store.Profile{Vars: make(map[string]string, len(srcProfile.Vars))}
		for k, v := range srcProfile.Vars {
			newProfile.Vars[k] = v
		}
		s.Profiles[dst] = newProfile
		if err := s.Save(); err != nil {
			fatal(err.Error())
		}
		printOK(fmt.Sprintf("Copied %q -> %q", src, dst))
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
