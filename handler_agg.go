package main

import (
	"context"
	"database/sql"
	"fmt"
	"gatorRSS/internal/database"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Command usage: agg <time> -> eg. 1h, 12m, 42s")
	}

	//Time Between Requests:
	TBR, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error parsing duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v:\n\n", TBR)

	ticker := time.NewTicker(TBR)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	newFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID: newFetch.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Error making feed: %w", err)
	}

	feed, err := fetchFeed(ctx, newFetch.Url)
	if err != nil {
		return err
	}

	fmt.Printf("%s - %s\n", feed.Channel.Title, feed.Channel.Description)

	for _, item := range feed.Channel.Item {

		var publishedAt sql.NullTime
		pubDate, err := parseTime(item.PubDate)
		publishedAt.Time = pubDate
		if err != nil {
			publishedAt.Valid = false
		} else {
			publishedAt.Valid = true
		}

		post, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID: uuid.NullUUID{
				UUID:  newFetch.ID,
				Valid: true,
			},
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			return fmt.Errorf("Error creating post: %w\n", err)
		} else {
			fmt.Printf("- %s\n", post.Title)
		}
	}
	fmt.Println("")
	return nil
}

func parseTime(pubTime string) (time.Time, error) {
	var layout = time.RFC1123
	if strings.HasSuffix(pubTime, "00") {
		layout = time.RFC1123Z
	}

	parsedTime, err := time.Parse(layout, pubTime)
	if err != nil {
		fmt.Printf("Error parsing date. Format: %s\n", pubTime)
		return time.Time{}, err
	}

	return parsedTime, nil
}
