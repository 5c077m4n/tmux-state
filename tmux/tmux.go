// Package tmux utils and types
package tmux

import (
	"bytes"
	"os/exec"
)

const Delimiter = "\u241f"
const listFormat = "#{session_id}" + Delimiter +
	"#{session_name}" + Delimiter +
	"#{session_attached}" + Delimiter +
	"#{window_id}" + Delimiter +
	"#{window_index}" + Delimiter +
	"#{window_name}" + Delimiter +
	"#{window_active}" + Delimiter +
	"#{window_layout}" + Delimiter +
	"#{pane_id}" + Delimiter +
	"#{pane_index}" + Delimiter +
	"#{pane_title}" + Delimiter +
	"#{pane_current_command}" + Delimiter +
	"#{pane_current_path}" + Delimiter +
	"#{pane_active}"

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
