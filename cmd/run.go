package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/hunntr/w3env/internal/store"
)

var runCmd = &cobra.Command{
	Use:   "run <command> [args...]",
	Short: "Run a command with the active profile's variables injected",
	Long: `Run a command with all profile variables set as environment variables.
Useful in scripts where you don't want to eval exports:

  w3env run cast call $TARGET "owner()(address)"
  w3env run forge script script/Exploit.s.sol --broadcast`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
			_ = cmd.Help()
			return
		}
		s, err := store.Load()
		if err != nil {
			fatal(err.Error())
		}
		p, err := s.ActiveProfile()
		if err != nil {
			fatal(err.Error())
		}
		env := os.Environ()
		for k, v := range p.Vars {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		env = append(env, fmt.Sprintf("W3ENV_ACTIVE=%s", s.Active))

		c := exec.Command(args[0], args[1:]...)
		c.Env = env
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
