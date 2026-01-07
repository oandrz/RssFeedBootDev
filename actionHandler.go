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

func handlerAddFeed(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("feed name and url argument is required")
	}

	user, err := state.dbQueriesData.GetUser(context.Background(), state.configData.CurrentUser)
	if err != nil {
		return fmt.Errorf("error on handler add feed get user: %v", err)
	}

	feedName := cmd.args[0]
	feedUrl := cmd.args[1]

	feed, err := state.dbQueriesData.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   feedName,
		Url:    feedUrl,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error on handler add feed: %v", err)
	}

	err = handlerFollow(state, command{command: "follow", args: []string{feedUrl}}, user)
	if err != nil {
		return fmt.Errorf("error on handler add feed when follow: %v", err)
	}

	fmt.Printf("Feed %s and url: %s successfully added\n", feed.Name, feed.Url)

	return nil
}

func handlerFeeds(state *state, cmd command) error {
	feeds, err := state.dbQueriesData.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error on handler feeds: %v", err)
	}
	for _, feed := range feeds {
		user, err := state.dbQueriesData.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error on handler feeds: %v", err)
		}
		fmt.Printf("Feed Title: %s\nFeed Url: %s\n", feed.Name, feed.Url)
		fmt.Printf("Created by %s\n", user.Name)
	}

	return nil
}

func handlerFollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("url argument is required")
	}

	user, err := state.dbQueriesData.GetUser(context.Background(), state.configData.CurrentUser)
	if err != nil {
		return fmt.Errorf("error on handler follow get user: %v", err)
	}

	feed, err := state.dbQueriesData.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error on handler follow get feed by url: %v", err)
	}

	feedFollow, err := state.dbQueriesData.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error on handler follow create feed follow: %v", err)
	}

	fmt.Printf("Feed %s successfully followed by %s\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFollowing(state *state, cmd command, user database.User) error {
	user, err := state.dbQueriesData.GetUser(context.Background(), state.configData.CurrentUser)
	if err != nil {
		return fmt.Errorf("error on handler following get user: %v", err)
	}

	feedFollows, err := state.dbQueriesData.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error on handler following get feed follows for user: %v", err)
	}

	fmt.Println("Feed followed by user:")
	for idx, feedFollow := range feedFollows {
		fmt.Printf("%d Feed Title: %s\nFeed Url: %s\n", idx+1, feedFollow.FeedName, feedFollow.FeedUrl)
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.dbQueriesData.GetUser(context.Background(), s.configData.CurrentUser)
		if err != nil {
			return fmt.Errorf("error on handler following get user: %v", err)
		}

		return handler(s, c, user)
	}
}
