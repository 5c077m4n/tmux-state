package main

import (
	"github.com/5c077m4n/tmux-state/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
