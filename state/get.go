// Package state getter/setter
package state

import (
	"errors"

	"github.com/5c077m4n/tmux-state/tmux"
)

var ErrRunTmux = errors.New("failed to run Tmux")

func GetTmuxState() (tmux.State, error) {
	stdout, _, err := tmux.ListPanes()
	if err != nil {
		return nil, errors.Join(ErrRunTmux, err)
	}

	return tmux.StateFrom(stdout)
}

