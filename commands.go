package main

import "fmt"

type command struct {
	name string
	args []string
}

type commandHandler func(*state, command) error

func loginHandler(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("error login handler expects a single argument, the username")
	}
	newUserName := cmd.args[0]
	err := s.cfg.SetUser(newUserName)
	if err != nil {
		return err
	}
	fmt.Printf("Current user has been set to %s\n", s.cfg.CurrentUser)
	return nil
}

type commands struct {
	commandHandlers map[string]commandHandler
}

func NewCommands() commands {

	return commands{
		commandHandlers: make(map[string]commandHandler),
	}
}

func (c commands) run(s *state, cmd command) error {
	commandHandler, ok := c.commandHandlers[cmd.name]
	if !ok {
		return fmt.Errorf("error no handler for [%s] command", cmd.name)
	}

	err := commandHandler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c commands) register(name string, f commandHandler) error {
	c.commandHandlers[name] = f
	return nil
}
