package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oleshko-g/gatorcli/internal/database"
)

const defaultRSSURL = "https://www.wagslane.dev/index.xml"

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
	ctx := context.Background()

	existingUser, errGetUserByName := s.db.GetUserByName(ctx, newUserName)

	userExists := existingUser != database.User{} &&
		errGetUserByName != sql.ErrNoRows
	if !userExists {
		return fmt.Errorf("error user doesn't exist")
	} else if errGetUserByName != nil {
		return errGetUserByName // error trying to check if user exists
	}

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

func resetUsersHandler(s *state, cmd command) error {
	argumentsNum := len(cmd.args)
	if argumentsNum > 0 {
		return fmt.Errorf("error reset handler expects no argument. Number of arguments %d", argumentsNum)
	}
	ctx := context.Background()

	errReset := s.db.ResetUsers(ctx)
	if errReset != nil {
		return errReset
	}
	fmt.Println("Users Table have been reset successfuly")
	return nil
}

func getUsersHandler(s *state, cmd command) error {
	argumentsNum := len(cmd.args)
	if argumentsNum > 0 {
		return fmt.Errorf("error handler expects no argument. Number of arguments %d", argumentsNum)
	}

	ctx := context.Background()

	users, errGetUsers := s.db.GetUsers(ctx)
	if errGetUsers != nil {
		return errGetUsers
	}

	for _, user := range users {
		current := ""
		if user.Name == s.cfg.CurrentUser {
			current = " (current)"
		}

		fmt.Printf("* %s%s\n", user.Name, current)
	}

	return nil
}

func aggHandler(s *state, cmd command) error {
	rssURL := defaultRSSURL

	argumentsNum := len(cmd.args)
	if argumentsNum == 1 {
		rssURL = cmd.args[0]
	}

	ctx := context.Background()

	rss, errFetchFeed := fetchFeed(ctx, rssURL)
	if errFetchFeed != nil {
		return errFetchFeed
	}

	fmt.Printf("%#v", *rss)

	return nil
}

func addfeedHandler(s *state, cmd command) error {
	argumentsNum := len(cmd.args)
	if argumentsNum != 2 {
		return fmt.Errorf("error addfeed handler expects 2 arguments. Number of arguments %d", argumentsNum)
	}

	ctx := context.Background()
	currentDBUser, errGetUserByName := s.db.GetUserByName(ctx, s.cfg.CurrentUser)
	if errGetUserByName != nil {
		return errGetUserByName
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	createdFeed, errCreateFeed := s.db.CreateFeed(ctx,
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      feedName,
			Url:       feedURL,
			UserID: uuid.NullUUID{
				UUID:  currentDBUser.ID,
				Valid: true,
			},
		})
	if errCreateFeed != nil {
		return errCreateFeed
	}

	fmt.Printf("%#v", createdFeed)

	return nil
}

type commands struct {
	commandHandlers map[string]commandHandler
}

func NewCommands() commands {
	cmds := commands{commandHandlers: make(map[string]commandHandler)}

	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetUsersHandler)
	cmds.register("users", getUsersHandler)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", addfeedHandler)

	return cmds
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
