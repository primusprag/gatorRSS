package main

import "errors"

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, exists := c.cmds[cmd.name]
	if !exists {
		return errors.New("Error: command not found")
	}

	return command(s, cmd)
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.cmds[name] = f
}
