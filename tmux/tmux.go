// Package tmux utils and types
package tmux

import (
	"bytes"
	"os/exec"
	"strings"
)

const Delimiter = "\u241f"

var listFormat = strings.Join(
	[]string{
		"#{session_id}",
		"#{session_name}",
		"#{session_attached}",
		"#{window_id}",
		"#{window_index}",
		"#{window_name}",
		"#{window_active}",
		"#{window_layout}",
		"#{pane_id}",
		"#{pane_index}",
		"#{pane_title}",
		"#{pane_current_command}",
		"#{pane_current_path}",
		"#{pane_active}",
	},
	Delimiter,
)

func Cmd(args ...string) (string, string, error) {
	cmd := exec.Command("tmux", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", "", err
	}
	return stdout.String(), stderr.String(), nil
}

func ListPanes() (string, string, error) {
	return Cmd("list-panes", "-a", "-F", listFormat)
}

func NewSession(sessionName, windowID string) (string, string, error) {
	return Cmd("new-session", "-d", "-s", sessionName, "-n", windowID)
}

func NewWindow(target, windowID string) (string, string, error) {
	return Cmd("new-window", "-t", target, "-n", windowID)
}

func SplitWindow(target string) (string, string, error) {
	return Cmd("split-window", "-t", target)
}

func SelectLayout(target, layout string) (string, string, error) {
	return Cmd("select-layout", "-t", target, layout)
}

func SendKeys(target, keys string) (string, string, error) {
	return Cmd("send-keys", "-t", target, keys)
}

func SelectPane(target string) (string, string, error) {
	return Cmd("select-pane", "-t", target)
}

func SelectWindow(target string) (string, string, error) {
	return Cmd("select-window", "-t", target)
}
