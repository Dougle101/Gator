package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Dougle101/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: 'gator login <username>', received %d arguments (expected 1)", len(cmd.args))
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error finding user, check that it exists: %v", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to %s", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: 'gator register <username>', received %d arguments (expecting 1)", len(cmd.args))
	}

	user_params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]}
	user, err := s.db.CreateUser(context.Background(), user_params)
	if err != nil {
		return fmt.Errorf("Error creating user: %v", err)
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("Created user: %s, ID: %s", user.Name, user.ID)
	return nil
}
