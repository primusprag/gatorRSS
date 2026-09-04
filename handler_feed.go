package main

import (
	"context"
	"fmt"
	"gatorRSS/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("command usage: addFeed <'name'> <'url'>")
	}

	ctx := context.Background()

	currentUserID, err := s.db.GetUserID(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error retrieving userID: %w", err)
	}

	feed, err := s.db.AddFeed(ctx, database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currentUserID,
	})
	if err != nil {
		return fmt.Errorf("Error adding feed: %w", err)
	}

	fmt.Printf(
		feed.ID.String(),
		feed.CreatedAt,
		feed.UpdatedAt,
		feed.Name,
		feed.Url,
		feed.UserID)

	return nil
}
