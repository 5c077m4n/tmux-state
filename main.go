package main

import (
	"encoding/json"
	"fmt"

	"github.com/5c077m4n/tmux-state/state"
)

func main() {
	s, err := state.GetTmuxState()
	if err != nil {
		panic(err)
	}
	stateJSON, _ := json.MarshalIndent(s, "", "\t")
	fmt.Printf("tmux state: %+v\n", string(stateJSON))
}
