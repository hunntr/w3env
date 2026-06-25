package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var installRC string

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install shell integration into your rc file (one-time setup)",
	Run: func(cmd *cobra.Command, args []string) {
		exe, err := os.Executable()
		if err != nil {
			fatal(err.Error())
		}
		exe, _ = filepath.EvalSymlinks(exe)

		rc := installRC
		if rc == "" {
			rc = detectRC()
		}
		if rc == "" {
			fatal("could not detect shell rc file - use --rc ~/.your_rc")
		}

		if containsW3env(rc) {
			printOK(fmt.Sprintf("already installed in %s", rc))
			printInfo(fmt.Sprintf("  run: source %s", rc))
			return
		}

		f, err := os.OpenFile(rc, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fatal(err.Error())
		}
		defer f.Close()

		fmt.Fprintf(f, "\n# w3env - https://github.com/hunntr/w3env\nsource <(%q init)\n", exe)

		printOK(fmt.Sprintf("installed in %s", rc))
		printInfo(fmt.Sprintf("  run: source %s", rc))
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove shell integration from your rc file",
	Run: func(cmd *cobra.Command, args []string) {
		rc := installRC
		if rc == "" {
			rc = detectRC()
		}
		if rc == "" {
			fatal("could not detect shell rc file - use --rc ~/.your_rc")
		}

		if !containsW3env(rc) {
			printInfo(fmt.Sprintf("w3env not found in %s", rc))
			return
		}

		if err := removeW3envLines(rc); err != nil {
			fatal(err.Error())
		}
		printOK(fmt.Sprintf("removed from %s", rc))
		printInfo(fmt.Sprintf("  run: source %s", rc))
	},
}

func detectRC() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	shell := filepath.Base(os.Getenv("SHELL"))
	switch shell {
	case "zsh":
		return filepath.Join(home, ".zshrc")
	case "bash":
		if p := filepath.Join(home, ".bashrc"); fileExists(p) {
			return p
		}
		return filepath.Join(home, ".bash_profile")
	default:
		return ""
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func containsW3env(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "w3env") && strings.Contains(scanner.Text(), "init") {
			return true
		}
	}
	return scanner.Err() == nil
}

func removeW3envLines(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "w3env") {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		f.Close()
		return err
	}
	f.Close()

	trimmed := strings.TrimRight(strings.Join(lines, "\n"), "\n")
	return os.WriteFile(path, []byte(trimmed+"\n"), 0644)
}

func init() {
	installCmd.Flags().StringVar(&installRC, "rc", "", "Path to rc file (default: auto-detect)")
	uninstallCmd.Flags().StringVar(&installRC, "rc", "", "Path to rc file (default: auto-detect)")
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
}
