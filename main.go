package main

import (
	"database/sql"
	"fmt"
	"gatorRSS/internal/config"
	"gatorRSS/internal/database"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	userCfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sql.Open("postgres", userCfg.DbUrl)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &userCfg,
	}

	cmds := commands{
		cmds: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	if len(os.Args) < 2 {
		fmt.Println("Error: command requires more than one argument.")
		os.Exit(1)
	}

	err = cmds.run(programState, command{name: os.Args[1], args: os.Args[2:]})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
