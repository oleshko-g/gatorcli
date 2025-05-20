package main

import (
	"fmt"
	"os"

	"github.com/oleshko-g/gatorcli/internal/config"
)

type state struct {
	cfg *config.Config
}

func setState() state {
	var s state
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("error reading config"))
		os.Exit(1)
	}
	s.cfg = &cfg
	return s
}

func parseArgs() command {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("error not enough arguments were provided."))
		os.Exit(1)
	}
	return command{
		name: os.Args[1],
		args: os.Args[2:],
	}
}

func main() {
	state := setState()
	cmds := NewCommands()
	cmds.register("login", loginHandler)
	cmd := parseArgs()
	err := cmds.run(&state, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
