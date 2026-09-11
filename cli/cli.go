// Package cli for CLI commands
package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/5c077m4n/tmux-state/state"
	"github.com/5c077m4n/tmux-state/tmux"
)

var (
	ErrStateFile    = errors.New("no state file found")
	ErrNoCommand    = errors.New("no command provided")
	ErrStateSave    = errors.New("could not save the tmux state")
	ErrStateRestore = errors.New("could not restore the tmux state")
)

func stateFilePath() (string, error) {
	dir, found := os.LookupEnv("XDG_STATE_HOME")
	if !found || dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", errors.Join(ErrStateFile, err)
		}

		dir = filepath.Join(home, ".local", "state")
	}

	return filepath.Join(dir, "tmux-state", "current.json"), nil
}

func saveState() error {
	s, err := state.GetTmuxState()
	if err != nil {
		return errors.Join(ErrStateSave, err)
	}

	path, err := stateFilePath()
	if err != nil {
		return errors.Join(ErrStateSave, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return errors.Join(ErrStateSave, err)
	}

	data, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return errors.Join(ErrStateSave, err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return errors.Join(ErrStateSave, err)
	}

	fmt.Printf("saved tmux state to %s\n", path)
	return nil
}

func restoreState() error {
	path, err := stateFilePath()
	if err != nil {
		return errors.Join(ErrStateRestore, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Join(ErrStateRestore, err)
	}

	tmuxState := &tmux.State{}
	if err := json.Unmarshal(data, tmuxState); err != nil {
		return errors.Join(ErrStateRestore, err)
	}

	if err := state.RestoreTmuxState(tmuxState); err != nil {
		return errors.Join(ErrStateRestore, err)
	}

	fmt.Printf("restored tmux state from %s\n", path)
	return nil
}

func Execute() error {
	save := flag.Bool("save", false, "save the current tmux state")
	restore := flag.Bool("restore", false, "restore the latest tmux state")
	flag.Parse()

	switch {
	case save != nil && *save:
		return saveState()
	case restore != nil && *restore:
		return restoreState()
	default:
		flag.Usage()
		return ErrNoCommand
	}
}
