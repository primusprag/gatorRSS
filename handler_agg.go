package main

import (
	"context"
	"database/sql"
	"fmt"
	"gatorRSS/internal/database"
	"time"
)

const wagsURL = "https://www.wagslane.dev/index.xml"

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Command usage: agg <time> -> eg. 1h, 12m, 42s")
	}

	//Time Between Requests:
	TBR, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Errro parsing duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v:\n\n", TBR)

	ticker := time.NewTicker(TBR)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID: nextFetch.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Error making feed: %w", err)
	}

	feed, err := fetchFeed(ctx, nextFetch.Url)
	if err != nil {
		return err
	}

	fmt.Printf("%s:\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Println(item.Title)
	}

	return nil
}
