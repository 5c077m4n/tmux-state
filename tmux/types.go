package tmux

import (
	"bytes"
	"errors"
	"log/slog"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/5c077m4n/tmux-state/layout"
)

type Pane struct {
	ID      string
	Index   int
	Title   string
	Command string
	Path    string
	Active  bool
}
type Window struct {
	ID     string
	Index  int
	Name   string
	Active bool
	Layout *layout.Layout
	Panes  []Pane
}
type Session struct {
	ID      string
	Name    string
	Active  bool
	Windows []Window
}
type State []*Session

var ErrParseState = errors.New("could not parse the Tmux state")

func getProcFromTTY(tty string) (string, error) {
	cmd := exec.Command("/bin/ps", "-t", tty, "-o", "command")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errMsg := stderr.String(); errMsg != "" {
			return "", errors.Join(errors.New(errMsg), err)
		}
		return "", err
	}
	if errMsg := stderr.String(); errMsg != "" {
		return "", errors.New(errMsg)
	}

	procs := slices.Collect(strings.Lines(stdout.String()))
	leafProc := procs[len(procs)-1]
	return strings.TrimSpace(leafProc), nil
}

func StateFrom(stdout string) (State, error) {
	sessionMap := map[string]*Session{}
	windowMap := map[string]*Window{}
	sessions := State{}

	for line := range strings.Lines(strings.TrimSpace(stdout)) {
		if line == "" {
			continue
		}

		parts := strings.Split(line, Delimiter)
		if len(parts) < 14 {
			continue
		}

		sID := parts[0]
		sName := parts[1]
		sAttached := parts[2] == "1"
		wID := parts[3]
		wIdx, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, errors.Join(ErrParseState, err)
		}
		wName := parts[5]
		wActive := parts[6] == "1"
		wLayout := parts[7]
		pID := parts[8]
		pIdx, err := strconv.Atoi(parts[9])
		if err != nil {
			return nil, errors.Join(ErrParseState, err)
		}
		pTitle := parts[10]
		pCmd := parts[11]
		pTTY := parts[12]
		pPath := parts[13]
		pActive := parts[14] == "1"

		sess, ok := sessionMap[sID]
		if !ok {
			sess = &Session{
				ID:     sID,
				Name:   sName,
				Active: sAttached,
			}
			sessionMap[sID] = sess
			sessions = append(sessions, sess)
		}

		win, ok := windowMap[wID]
		if !ok {
			parsedLayout, err := layout.Parse(wLayout)
			if err != nil {
				slog.Warn("could not parse layout", slog.Any("error", err))
			}

			sess.Windows = append(sess.Windows, Window{
				ID:     wID,
				Index:  wIdx,
				Name:   wName,
				Active: wActive,
				Layout: parsedLayout,
			})
			windowMap[wID] = &sess.Windows[len(sess.Windows)-1]
			win = windowMap[wID]
		}

		proc, err := getProcFromTTY(pTTY)
		if err != nil {
			slog.Warn(
				"could not get process",
				slog.String("tty", pTTY),
				slog.Any("error", err),
			)
			proc = pCmd
		}
		win.Panes = append(win.Panes, Pane{
			ID:      pID,
			Index:   pIdx,
			Title:   pTitle,
			Command: proc,
			Path:    pPath,
			Active:  pActive,
		})
	}

	return sessions, nil
}
