package main

import (
	"bootDevGoRss/internal/config"
	"errors"
	"fmt"
)

type state struct {
	configData *config.Config
}

/*
For example, in the case of the login command,
the name would be "login" and the handler will expect the arguments slice to contain one string, the username.
*/
type command struct {
	command string
	args    []string
}

type commands struct {
	// This will be a map of command names to their handler functions.
	mapper map[string]func(*state, command) error
}

/*
This will be the function signature of all command handlers
*/
func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username argument is required")
	}

	if err := state.configData.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Println("Logged in as", cmd.args[0])
	return nil
}

/*
*
This method runs a given command with the provided state if it exists.
*/
func (c *commands) run(state *state, cmd command) error {
	handler, ok := c.mapper[cmd.command]
	if !ok {
		return fmt.Errorf("command %s does not exist", cmd.command)
	}
	return handler(state, cmd)
}

/*
*
This method registers a new handler function for a command name.
*/
func (c *commands) register(name string, f func(*state, command) error) error {
	_, ok := c.mapper[name]
	if ok {
		return fmt.Errorf("command %s already exists", name)
	}

	c.mapper[name] = f
	return nil
}
