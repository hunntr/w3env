package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "w3env",
	Short: "Web3 environment manager for pentesters and CTF players",
	Long: `w3env - manage named profiles of Web3 variables (RPC URLs, contract addresses, private keys, chain IDs...) and switch between them instantly.

Quick start:
  w3env new htb-challenge1
  w3env use htb-challenge1
  w3env set RPC_URL   http://rpc.target:8545
  w3env set TARGET    0xDeAdBeEf...
  w3env set PRIVATE_KEY 0x1234...
  w3env show

Foundry/cast integration - w3env var names map directly:
  RPC_URL      -> ETH_RPC_URL    (cast / forge)
  PRIVATE_KEY  -> picked up by cast send --private-key

Shell integration (auto-exports vars on profile switch):
  echo 'source <(w3env init)' >> ~/.zshrc  && source ~/.zshrc`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
