// Package tmux utils and types
package tmux

import (
	"bytes"
	"errors"
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

func runTMuxCommand(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}
	if errMsg := stderr.String(); errMsg != "" {
		return "", errors.New(errMsg)
	}
	return stdout.String(), nil
}

func ListPanes() (string, error) {
	return runTMuxCommand("list-panes", "-a", "-F", listFormat)
}

func NewSession(sessionName, windowID string) (string, error) {
	return runTMuxCommand("new-session", "-d", "-s", sessionName, "-n", windowID)
}

func NewWindow(target, windowID string) (string, error) {
	return runTMuxCommand("new-window", "-t", target, "-n", windowID)
}

func SplitWindow(target string) (string, error) {
	return runTMuxCommand("split-window", "-t", target)
}

func SelectLayout(target, layout string) (string, error) {
	return runTMuxCommand("select-layout", "-t", target, layout)
}

func SendKeys(target, keys string) (string, error) {
	return runTMuxCommand("send-keys", "-t", target, keys)
}

func SelectPane(target string) (string, error) {
	return runTMuxCommand("select-pane", "-t", target)
}

func SelectWindow(target string) (string, error) {
	return runTMuxCommand("select-window", "-t", target)
}
