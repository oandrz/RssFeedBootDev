package main

import (
	"bootDevGoRss/internal/database"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username argument is required")
	}

	_, err := state.dbQueriesData.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user %s does not exists", cmd.args[0])
	}

	if err := state.configData.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Println("Logged in as", cmd.args[0])
	return nil
}

func handlerRegister(state *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username argument is required")
	}

	_, err := state.dbQueriesData.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		return fmt.Errorf("user %s already exists", cmd.args[0])
	}

	user, err := state.dbQueriesData.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("cannot create user: %v", err)
	}

	if err := state.configData.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Printf("User success fully created and registered:\n id: %v, name: %v\n", user.ID, user.Name)
	return nil
}

func handlerDelete(state *state, cmd command) error {
	err := state.dbQueriesData.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handlerGetUsers(state *state, cmd command) error {
	users, err := state.dbQueriesData.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == state.configData.CurrentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

const hardCodedUrl = "https://www.wagslane.dev/index.xml"

func handlerAggCommand(state *state, cmd command) error {
	feeds, err := fetchFeed(context.Background(), hardCodedUrl)
	if err != nil {
		return err
	}
	fmt.Println(feeds)
	return nil
}
