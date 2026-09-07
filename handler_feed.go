package main

import (
	"context"
	"fmt"
	"gatorRSS/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Command usage: addFeed <'name'> <'url'>")
	}

	ctx := context.Background()

	feed, err := s.db.AddFeed(ctx, database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error adding feed: %w", err)
	}

	err = handlerFollow(s, command{
		name: "follow",
		args: []string{feed.Url},
	}, user)

	fmt.Printf(`Added feed:
	ID: %s
	Created At: %v
	Updated at: %v
	Name: %s
	URL: %s
	User ID: %v`,
		feed.ID.String(),
		feed.CreatedAt,
		feed.UpdatedAt,
		feed.Name,
		feed.Url,
		feed.UserID)

	return nil
}

func handlerFeeds(s *state, _ command) error {
	ctx := context.Background()

	feeds, err := s.db.ListFeeds(ctx)
	if err != nil {
		return fmt.Errorf("Error retreving feeds: %w", err)
	}

	for _, feed := range feeds {
		feedUser, err := s.db.GetNameByID(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("Error retreving username: %w", err)
		}
		fmt.Printf("Name: %s URL: %s User: %s\n", feed.Name, feed.Url, feedUser)
	}

	return nil
}
