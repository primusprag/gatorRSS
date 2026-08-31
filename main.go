package main

import (
	"fmt"
	"gatorRSS/internal/config"
	"os"
)

type state struct {
	cfg *config.Config
}

func main() {

	userCfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	programState := &state{
		cfg: &userCfg,
	}

	cmds := commands{
		cmds: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Error: command requires more than one argument.")
		os.Exit(1)
	}

	command := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(programState, command)
	if err != nil {
		fmt.Printf("Error executing: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
