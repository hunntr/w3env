package cmd

import (
	"fmt"
	"os"
)

var noColor = os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb"

const (
	cReset  = "\033[0m"
	cBold   = "\033[1m"
	cRed    = "\033[31m"
	cGreen  = "\033[32m"
	cYellow = "\033[33m"
	cCyan   = "\033[36m"
	cGray   = "\033[90m"
)

func col(c, s string) string {
	if noColor {
		return s
	}
	return c + s + cReset
}

func printOK(msg string)   { fmt.Println(col(cGreen, "✓ "+msg)) }
func printInfo(msg string) { fmt.Println(col(cCyan, msg)) }
func printWarn(msg string) { fmt.Fprintln(os.Stderr, col(cYellow, "⚠ "+msg)) }
func fatal(msg string) {
	fmt.Fprintln(os.Stderr, col(cRed, "✗ "+msg))
	os.Exit(1)
}
