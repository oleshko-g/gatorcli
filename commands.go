package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
		return fmt.Errorf("error register handler expects a single argument, the username")
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
	argumentsNum := len(cmd.args)
	if argumentsNum != 1 {
		return fmt.Errorf("error 'agg' handler expects 1 arguments. Number of arguments %d", argumentsNum)
	}

	arg := cmd.args[0]

	time_between_reqs, errParseDuration := time.ParseDuration(arg)
	if errParseDuration != nil {
		return errParseDuration
	}
	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func addfeedHandler(s *state, cmd command, u database.User) error {
	argumentsNum := len(cmd.args)
	if argumentsNum != 2 {
		return fmt.Errorf("error addfeed handler expects 2 arguments. Number of arguments %d", argumentsNum)
	}

	ctx := context.Background()

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
				UUID:  u.ID,
				Valid: true,
			},
		})

	followRecord, errCreateFeedFollow := s.db.CreateFeedFollow(ctx,
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    u.ID,
			FeedID:    createdFeed.ID,
		})
	if errCreateFeedFollow != nil {
		return errCreateFeedFollow
	}

	printCreatedFeedFollow(followRecord)

	if errCreateFeed != nil {
		return errCreateFeed
	}

	fmt.Printf("%+v", createdFeed)

	return nil
}

func printFeeds(f []database.GetFeedsUsersRow) {
	if len(f) == 0 {
		return
	}
	for _, v := range f {
		fmt.Printf("Feed Name: %s. Feed URL: %s. User Name: %s\n",
			v.Feed.Name,
			v.Feed.Url,
			v.User.Name,
		)
	}
}

func feedsHandler(s *state, cmd command) error {
	argumentsNum := len(cmd.args)
	if argumentsNum != 0 {
		return fmt.Errorf("error 'feeds' handler expects 0 arguments. Number of arguments %d", argumentsNum)
	}
	ctx := context.Background()

	feeds, errGetFeedsUsers := s.db.GetFeedsUsers(ctx)
	if errGetFeedsUsers != nil {
		return errGetFeedsUsers
	}
	printFeeds(feeds)

	return nil
}

func printCreatedFeedFollow(f database.CreateFeedFollowRow) {
	fmt.Printf("Feed name: %s. User Name: %s\n", f.Feed.Name, f.User.Name)
}

func followHandler(s *state, cmd command, u database.User) error {
	argumentsNum := len(cmd.args)
	if argumentsNum != 1 {
		return fmt.Errorf("error 'follow' handler expects 1 argument. Number of input arguments %d", argumentsNum)
	}
	ctx := context.Background()

	feedURL := cmd.args[0]
	feedData, errGetFeedByURL := s.db.GetFeedByURL(ctx, feedURL)
	if errGetFeedByURL != nil {
		return errGetFeedByURL
	}

	followRecord, errCreateFeedFollow := s.db.CreateFeedFollow(ctx,
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    u.ID,
			FeedID:    feedData.ID,
		})
	if errCreateFeedFollow != nil {
		return errCreateFeedFollow
	}

	printCreatedFeedFollow(followRecord)

	return nil
}

func unfollowHandler(s *state, cmd command, u database.User) error {
	argumentsNum := len(cmd.args)
	if argumentsNum != 1 {
		return fmt.Errorf("error 'unfollow' handler expects 1 argument. Number of input arguments %d", argumentsNum)
	}
	ctx := context.Background()

	feedURL := cmd.args[0]
	feedData, errGetFeedByURL := s.db.GetFeedByURL(ctx, feedURL)
	if errGetFeedByURL != nil {
		return errGetFeedByURL
	}

	_, errDel := s.db.DeleteFeedFollowUser(ctx,
		database.DeleteFeedFollowUserParams{
			UserID: u.ID,
			FeedID: feedData.ID,
		})
	if errDel != nil {
		return errDel
	}

	fmt.Printf("'%s' feed has been unfollowed", feedData.Name)

	return nil
}
func followingHandler(s *state, cmd command, u database.User) error {
	argumentsNum := len(cmd.args)
	if argumentsNum != 0 {
		return fmt.Errorf("error 'feeds' handler expects 0 arguments. Number of arguments %d", argumentsNum)
	}

	ctx := context.Background()

	feedFollows, errGetFeedFollowUser := s.db.GetFeedFollowUser(ctx, u.ID)
	if errGetFeedFollowUser != nil {
		return errGetFeedFollowUser
	}
	for _, feedFollow := range feedFollows {
		fmt.Printf("User Name: %s. Feed Name: %s\n", feedFollow.User.Name, feedFollow.Feed.Name)
	}

	return nil
}

func browseHandler(s *state, cmd command, u database.User) error {
	argumentsNum := len(cmd.args)
	if argumentsNum > 1 {
		return fmt.Errorf("error 'browse' handler expects 0 or 1 argument(s). Number of arguments %d", argumentsNum)
	}

	limit := 2
	if argumentsNum == 1 {
		if inputLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = inputLimit
		}
	}

	ctx := context.Background()

	posts, errGetPostsForUser := s.db.GetPostsForUser(ctx,
		database.GetPostsForUserParams{
			UserID: u.ID,
			Limit:  int32(limit),
		})

	if errGetPostsForUser != nil {
		return errGetPostsForUser
	}
	for i, feed := range posts {
		fmt.Printf("Post #%d:%s\n", i+1, feed.Description)
	}

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
	cmds.register("addfeed", middleWareLoggedIn(addfeedHandler))
	cmds.register("feeds", feedsHandler)
	cmds.register("follow", middleWareLoggedIn(followHandler))
	cmds.register("unfollow", middleWareLoggedIn(unfollowHandler))
	cmds.register("following", middleWareLoggedIn(followingHandler))
	cmds.register("browse", middleWareLoggedIn(browseHandler))

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
