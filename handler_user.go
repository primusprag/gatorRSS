package main

import (
	"context"
	"errors"
	"fmt"
	"gatorRSS/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Command usage: login <username>")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return errors.New("Error: Unable to login, user not in database.")
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error: Couldn't set user: %w", err)
	}

	fmt.Println("Username has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Command only accepts single-name inputs: register <name>")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("Couldn't create user: %w", err)
	}
	fmt.Printf("User %s was created\n", user.Name)

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}
	fmt.Printf("User %s was set\n", s.cfg.CurrentUserName)

	return nil
}

func handlerGetUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error retreving users: %w", err)
	}

	for _, name := range users {
		if name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", name)
		} else {
			fmt.Printf("* %s\n", name)
		}
	}

	return nil
}
