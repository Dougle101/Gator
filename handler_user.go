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
		return fmt.Errorf("Error finding user, check that it exists: %v\n", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to %s\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: 'gator register <username>', received %d arguments (expecting 1)\n", len(cmd.args))
	}

	user_params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]}
	user, err := s.db.CreateUser(context.Background(), user_params)
	if err != nil {
		return fmt.Errorf("Error creating user: %v\n", err)
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("Created user: %s, ID: %s\n", user.Name, user.ID)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: 'gator reset', received %d arguments (expecting 0)\n", len(cmd.args))
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting users: %v\n", err)
	}
	fmt.Println("All users purged from database.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: 'gator users', received %d arguments (expecting 0)\n", len(cmd.args))
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error fetching users: %v\n", err)
	}

	for _, name := range users {
		if s.cfg.CurrentUserName == name {
			fmt.Printf("%s (current)\n", name)
		} else {
			fmt.Printf("%s\n", name)
		}
	}

	return nil
}
