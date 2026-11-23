package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Dougle101/Gator/internal/config"
	"github.com/Dougle101/Gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	dbURL := cfg.DbUrl
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	s := &state{cfg: &cfg, db: database.New(db)}
	cmds := commands{cmd_map: make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

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
