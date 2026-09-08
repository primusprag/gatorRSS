package main

import (
	"context"
	"fmt"
	"gatorRSS/internal/database"
	"strconv"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int64
	if len(cmd.args) >= 1 {
		var err error
		limit, err = strconv.ParseInt(cmd.args[0], 8, 32)
		if err != nil {
			return fmt.Errorf("Command usage: browse (<length>)")
		}

	} else {
		limit = 2
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("%s - %s\n", post.FeedName.String, post.Title)
	}

	return nil
}

func handlerResetPosts(s *state, _ command) error {
	err := s.db.ResetPosts(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting posts: %w", err)
	}
	fmt.Println("Posts reset")
	return nil
}
