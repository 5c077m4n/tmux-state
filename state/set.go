package state

import (
	"fmt"
	"log/slog"

	"github.com/5c077m4n/tmux-state/tmux"
)

func RestoreTmuxState(state *tmux.State) error {
	for _, s := range *state {
		if len(s.Windows) == 0 {
			continue
		}

		firstWin := s.Windows[0]
		if _, err := tmux.NewSession(s.Name, firstWin.Name); err != nil {
			slog.Error(
				"could not create session",
				slog.Any("error", err),
				slog.String("session ID", s.ID),
			)
			continue
		}

		for windowIdx, w := range s.Windows {
			if windowIdx > 0 {
				if _, err := tmux.NewWindow(s.Name, w.Name); err != nil {
					slog.Error(
						"could not create window",
						slog.Any("error", err),
						slog.String("window ID", w.ID),
					)
					continue
				}
			}

			targetWindow := fmt.Sprintf("%s:%d", s.Name, w.Index)

			for pIdx := 1; pIdx < len(w.Panes); pIdx++ {
				if _, err := tmux.SplitWindow(targetWindow); err != nil {
					slog.Error(
						"could not split window",
						slog.Any("error", err),
						slog.String("window", targetWindow),
					)
				}
			}

			if w.Layout != nil {
				if _, err := tmux.SelectLayout(targetWindow, w.Layout.String()); err != nil {
					slog.Error(
						"could not select layout",
						slog.Any("error", err),
						slog.String("window", targetWindow),
					)
				}
			}

			for _, p := range w.Panes {
				targetPane := fmt.Sprintf("%s.%d", targetWindow, p.Index)

				if p.Path != "" {
					cdCmd := fmt.Sprintf("cd %q && clear\n", p.Path)
					if _, err := tmux.SendKeys(targetPane, cdCmd); err != nil {
						slog.Error(
							"could not send path",
							slog.Any("error", err),
							slog.String("pane", targetPane),
						)
					}
				}

				if p.Command != "bash" && p.Command != "zsh" && p.Command != "fish" &&
					p.Command != "" {
					if _, err := tmux.SendKeys(targetPane, p.Command+"\n"); err != nil {
						slog.Error(
							"could not send command",
							slog.Any("error", err),
							slog.String("pane", targetPane),
						)
					}
				}

				if p.Title != "" {
					if _, err := tmux.SelectPaneTitle(targetPane, p.Title); err != nil {
						slog.Error(
							"could not set pane title",
							slog.Any("error", err),
							slog.String("pane", targetPane),
						)
					}
				}

				if p.Active {
					if _, err := tmux.SelectPane(targetPane); err != nil {
						slog.Error(
							"could not select pane",
							slog.Any("error", err),
							slog.String("pane", targetPane),
						)
					}
				}
			}

			if w.Active {
				if _, err := tmux.SelectWindow(targetWindow); err != nil {
					slog.Error(
						"could not select window",
						slog.Any("error", err),
						slog.String("window", targetWindow),
					)
				}
			}
		}
	}
	return nil
}
