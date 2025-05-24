package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oleshko-g/gatorcli/internal/database"
)

type command struct {
	name string
	args []string
}

type commandHandler func(*state, command) error

func loginHandler(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("error login handler expects a single argument, the username")
	}
	newUserName := cmd.args[0]
	err := s.cfg.SetUser(newUserName)
	if err != nil {
		return err
	}
	fmt.Printf("Current user has been set to %s\n", s.cfg.CurrentUser)
	return nil
}

func registerHandler(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("error login handler expects a single argument, the username")
	}

	newUserName := cmd.args[0]
	ctx := context.Background()

	existingUser, errGetUserByName := s.db.GetUserByName(ctx, newUserName)
	userExists := existingUser != database.User{} &&
		errGetUserByName == nil
	if userExists {
		return fmt.Errorf("error user name already exists")
	}

	if errGetUserByName == sql.ErrNoRows {
		createdUser, errCreateUser := s.db.CreateUser(
			ctx,
			database.CreateUserParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      newUserName,
			},
		)
		if errCreateUser != nil {
			return errCreateUser
		}

		err := s.cfg.SetUser(createdUser.Name)
		if err != nil {
			return err
		}
		fmt.Printf("Current user has been set to %s\n", s.cfg.CurrentUser)

		return nil
	}

	return errGetUserByName // error trying to check if user exists

}

type commands struct {
	commandHandlers map[string]commandHandler
}

func NewCommands() commands {

	return commands{
		commandHandlers: make(map[string]commandHandler),
	}
}

func (c commands) run(s *state, cmd command) error {
	commandHandler, ok := c.commandHandlers[cmd.name]
	if !ok {
		return fmt.Errorf("error no handler for [%s] command", cmd.name)
	}

	err := commandHandler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c commands) register(name string, f commandHandler) error {
	c.commandHandlers[name] = f
	return nil
}
