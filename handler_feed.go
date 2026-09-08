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
User ID: %v
`,
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

func handlerArchiveFeed(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Command usage: deletefeed <feed_name>")
	}

	ctx := context.Background()

	feed, err := s.db.GetFeedByName(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving feed: %w", err)
	}

	archFeed, err := s.db.ArchiveFeed(ctx, database.ArchiveFeedParams{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: time.Now(),
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	})
	if err != nil {
		return fmt.Errorf("Error archiving feed: %w", err)
	}

	_, err = s.db.DeleteFeed(ctx, archFeed.ID)

	fmt.Printf("%s archived\n", archFeed.Name)
	return nil
}

func handlerArchiveFeeds(s *state, cmd command) error {
	for _, arg := range cmd.args {
		err := handlerArchiveFeed(s, command{
			name: "deletefeed",
			args: []string{arg},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func handlerRestoreFeed(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Command usage: deletefeed <feed_name>")
	}

	ctx := context.Background()

	feed, err := s.db.GetArchivedFeedByName(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving feed: %w", err)
	}

	resFeed, err := s.db.AddFeed(ctx, database.AddFeedParams{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: time.Now(),
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	})
	if err != nil {
		return fmt.Errorf("Error archiving feed: %w", err)
	}

	_, err = s.db.DeleteArchivedFeed(ctx, resFeed.ID)

	fmt.Printf("%s restored\n", resFeed.Name)
	return nil
}

func handlerRestoreFeeds(s *state, cmd command) error {
	for _, arg := range cmd.args {
		err := handlerRestoreFeed(s, command{
			name: "restorefeed",
			args: []string{arg},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
