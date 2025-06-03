package main

import (
	"context"

	"github.com/oleshko-g/gatorcli/internal/database"
)

type loggedInCommandHandler func(s *state, cmd command, u database.User) error

func middleWareLoggedIn(h loggedInCommandHandler) commandHandler {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		currentUserData, errGetUserByName := s.db.GetUserByName(ctx, s.cfg.CurrentUser)
		if errGetUserByName != nil {
			return errGetUserByName
		}

		errH := h(s, cmd, currentUserData)
		if errH != nil {
			return errH
		}
		return nil
	}
}
