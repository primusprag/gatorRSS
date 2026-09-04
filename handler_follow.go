package main

import (
	"context"
	"fmt"
	"gatorRSS/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Command usage: follow <url>")
	}
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %w", err)
	}

	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving feed: %w", err)
	}

	feed_follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating feed-follow table: %w", err)
	}

	fmt.Printf("User: %s - Feed: %s\n", feed_follow[0].UserName, feed_follow[0].FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Command takes no inputs")
	}

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %w", err)
	}

	following, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("Error retrieving list: %w", err)
	}

	for _, feed := range following {
		fmt.Printf("User: %s - Feed: %s\n", feed.UserName, feed.FeedName)
	}

	return nil
}
