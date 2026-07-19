package state

import (
	"fmt"

	"github.com/5c077m4n/tmux-state/tmux"
)

func RestoreTmuxState(state tmux.State) error {
	for _, s := range state {
		firstWin := s.Windows[0]
		_, _, err := tmux.NewSession(s.Name, firstWin.ID)
		if err != nil {
			return err
		}

		for windowIdx, w := range s.Windows {
			if windowIdx > 0 {
				_, _, err = tmux.NewWindow(s.Name, w.ID)
				if err != nil {
					return err
				}
			}

			targetWindow := fmt.Sprintf("%s:%d", s.Name, w.Index)

			for pIdx := 1; pIdx < len(w.Panes); pIdx++ {
				_, _, err = tmux.SplitWindow(targetWindow)
				if err != nil {
					return err
				}
			}

			if w.Layout != "" {
				_, _, err = tmux.SelectLayout(targetWindow, w.Layout)
				if err != nil {
					return err
				}
			}

			for _, p := range w.Panes {
				targetPane := fmt.Sprintf("%s.%d", targetWindow, p.Index)

				if p.Path != "" {
					cdCmd := fmt.Sprintf("cd %q && clear\n", p.Path)
					_, _, _ = tmux.SendKeys(targetPane, cdCmd)
				}

				if p.Command != "bash" && p.Command != "zsh" && p.Command != "fish" &&
					p.Command != "" {
					_, _, _ = tmux.SendKeys(targetPane, p.Command+"\n")
				}

				if p.Active {
					_, _, _ = tmux.SelectPane(targetPane)
				}
			}

			if w.Active {
				_, _, _ = tmux.SelectWindow(targetWindow)
			}
		}
	}
	return nil
}
