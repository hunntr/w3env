package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Print shell integration code to source",
	Long: `Print the shell function needed for automatic variable export on profile switch.

Setup (one time):
  echo 'source <(w3env init)' >> ~/.zshrc   # zsh
  echo 'source <(w3env init)' >> ~/.bashrc  # bash
  source ~/.zshrc`,
	Run: func(cmd *cobra.Command, args []string) {
		exe, err := os.Executable()
		if err != nil {
			exe = "w3env"
		} else {
			exe, _ = filepath.EvalSymlinks(exe)
		}

		fmt.Printf(`# w3env shell integration - https://github.com/hunntr/w3env
_W3ENV_BIN=%q

w3env() {
    "$_W3ENV_BIN" "$@"
}

_w3env_precmd() {
    local f="${HOME}/.config/w3env/pending-$$.sh"
    [[ -f "$f" ]] && { source "$f"; rm -f "$f"; }
}

_w3env_ps1() {
    [[ -z "$W3ENV_ACTIVE" ]] && return
    echo -n "[${W3ENV_ACTIVE}] "
    # echo -n "(${W3ENV_ACTIVE}) "
}

if [[ -n "$ZSH_VERSION" ]]; then
    setopt PROMPT_SUBST
    precmd_functions+=(_w3env_precmd)
    [[ -z "$_W3ENV_PROMPT_ORIG" ]] && _W3ENV_PROMPT_ORIG="$PROMPT"
    PROMPT='$(_w3env_ps1)'"$_W3ENV_PROMPT_ORIG"
elif [[ -n "$BASH_VERSION" ]]; then
    PROMPT_COMMAND="${PROMPT_COMMAND:+$PROMPT_COMMAND; }_w3env_precmd"
    [[ -z "$_W3ENV_PS1_ORIG" ]] && _W3ENV_PS1_ORIG="$PS1"
    PS1='$(_w3env_ps1)'"$_W3ENV_PS1_ORIG"
fi
`, exe)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
