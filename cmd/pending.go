package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hunntr/w3env/internal/store"
)

func writeActivation(s *store.State, newName string) {
	var sb strings.Builder
	prevName := os.Getenv("W3ENV_ACTIVE")
	if prevName != "" && prevName != newName {
		if prev, ok := s.Profiles[prevName]; ok {
			newVars := s.Profiles[newName].Vars
			for k := range prev.Vars {
				if _, kept := newVars[k]; !kept {
					fmt.Fprintf(&sb, "unset %s\n", k)
				}
			}
		}
	}
	p := s.Profiles[newName]
	for k, v := range p.Vars {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	fmt.Fprintf(&sb, "export W3ENV_ACTIVE=%q\n", newName)
	store.WritePending(os.Getppid(), sb.String())
}

func writeDeactivationVars(vars map[string]string) {
	var sb strings.Builder
	for k := range vars {
		fmt.Fprintf(&sb, "unset %s\n", k)
	}
	sb.WriteString("unset W3ENV_ACTIVE\n")
	store.WritePending(os.Getppid(), sb.String())
}

func writeDeactivation(s *store.State, name string) {
	var sb strings.Builder
	if p, ok := s.Profiles[name]; ok {
		for k := range p.Vars {
			fmt.Fprintf(&sb, "unset %s\n", k)
		}
	}
	sb.WriteString("unset W3ENV_ACTIVE\n")
	store.WritePending(os.Getppid(), sb.String())
}
