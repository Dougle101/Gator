package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Dougle101/Gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s := &state{cfg: &cfg}
	cmds := commands{cmd_map: make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)

	user_cmd := os.Args
	if len(user_cmd) < 2 {
		fmt.Printf("Fatal error: You must provide a command.")
		os.Exit(1)
	}

	err = cmds.run(s, command{name: user_cmd[1], args: user_cmd[2:]})
	if err != nil {
		fmt.Printf("Fatal error: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
