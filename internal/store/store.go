package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Profile struct {
	Vars map[string]string `json:"vars"`
}

type State struct {
	Active   string             `json:"active"`
	Profiles map[string]Profile `json:"profiles"`
}

var (
	ErrNoActiveProfile = errors.New("no active profile - run: w3env use <profile>")
	ErrProfileNotFound = errors.New("profile not found")
)

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "w3env")
	return dir, os.MkdirAll(dir, 0700)
}

func Load() (*State, error) {
	dir, err := ConfigDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "state.json"))
	if os.IsNotExist(err) {
		return &State{Profiles: make(map[string]Profile)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read state: %w", err)
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}
	if s.Profiles == nil {
		s.Profiles = make(map[string]Profile)
	}
	return &s, nil
}

func (s *State) Save() error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "state.json"), data, 0600)
}

func (s *State) ActiveProfile() (*Profile, error) {
	if s.Active == "" {
		return nil, ErrNoActiveProfile
	}
	p, ok := s.Profiles[s.Active]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrProfileNotFound, s.Active)
	}
	return &p, nil
}

func WritePending(ppid int, content string) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, fmt.Sprintf("pending-%d.sh", ppid))
	return os.WriteFile(path, []byte(content), 0600)
}

func (s *State) ProfileNames() []string {
	names := make([]string, 0, len(s.Profiles))
	for n := range s.Profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
