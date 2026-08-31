package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 || len(cmd.args) > 1 {
		return fmt.Errorf("Error: the command 'login' expects only username")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set user: %w", err)
	}

	fmt.Println("Username has been set.")
	return nil
}
