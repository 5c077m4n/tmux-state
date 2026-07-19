package tmux

import (
	"errors"
	"strconv"
	"strings"
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
	Name   bool
	Active bool
	Layout string
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

func StateFrom(stdout string) (State, error) {
	sessionMap := map[string]*Session{}
	windowMap := map[string]*Window{}
	sessions := State{}

	for line := range strings.SplitSeq(strings.TrimSpace(stdout), "\n") {
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
		wName := parts[5] == "1"
		wActive := parts[6] == "1"
		wLayout := parts[7]
		pID := parts[8]
		pIdx, err := strconv.Atoi(parts[9])
		if err != nil {
			return nil, errors.Join(ErrParseState, err)
		}
		pTitle := parts[10]
		pCmd := parts[11]
		pPath := parts[12]
		pActive := parts[13] == "1"

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
			sess.Windows = append(sess.Windows, Window{
				ID:     wID,
				Index:  wIdx,
				Name:   wName,
				Active: wActive,
				Layout: wLayout,
			})
			windowMap[wID] = &sess.Windows[len(sess.Windows)-1]
			win = windowMap[wID]
		}

		win.Panes = append(win.Panes, Pane{
			ID:      pID,
			Index:   pIdx,
			Title:   pTitle,
			Command: pCmd,
			Path:    pPath,
			Active:  pActive,
		})
	}

	return sessions, nil
}
