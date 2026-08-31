package main

import (
	"fmt"
	"gatorRSS/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = cfg.SetUser("d4l4-33"); err != nil {
		fmt.Println(err)
		return
	}
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("db_url : %s\nusername: %s\n", cfg.DbUrl, cfg.CurrentUserName)
}
