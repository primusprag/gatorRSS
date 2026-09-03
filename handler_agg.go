package main

import (
	"context"
	"fmt"
)

const wagsURL = "https://www.wagslane.dev/index.xml"

func handlerAgg(s *state, cmd command) error {
	RSS, err := fetchFeed(context.Background(), wagsURL)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %w", err)
	}

	fmt.Println(RSS)
	return nil
}
