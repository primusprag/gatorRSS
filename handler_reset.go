package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error reseting table: %w", err)
	}
	fmt.Println("Users reset")
	return nil
}
