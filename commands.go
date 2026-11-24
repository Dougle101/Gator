package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmd_map map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	fnc, ok := c.cmd_map[cmd.name]
	if !ok {
		return fmt.Errorf("Command '%s' not found\n", cmd.name)
	}
	err := fnc(s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd_map[name] = f
}
