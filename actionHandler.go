package main

import (
	"bootDevGoRss/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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
	if len(cmd.args) < 1 {
		return errors.New("agg command needs time between request")
	}

	timeParam := cmd.args[0]
	timeBetweenRequests, err := time.ParseDuration(timeParam)
	ticker := time.NewTicker(timeBetweenRequests)

	// channel, blocked until got the c channel emit the item
	for ; ; <-ticker.C {
		fmt.Printf("Collecting feeds every 1m0s")

		err = scrapeFeeds(state)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(state *state) error {
	nextFeed, err := state.dbQueriesData.GetNextFeedToFetched(context.Background())
	if err != nil {
		return fmt.Errorf("error when scrape feed when get next feed %v", err)
	}

	err = state.dbQueriesData.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true, // true means the value is not NULL
		},
		Url: nextFeed.Url,
	})
	if err != nil {
		return fmt.Errorf("error when scrape feed on mark feed fetched %v", err)
	}

	feeds, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error when scrape feed on fetch feed fetched %v", err)
	}

	fmt.Printf("Title for feed %s\n", feeds.Channel.Title)
	for idx, item := range feeds.Channel.Item {
		fmt.Printf("title: %s\n", item.Title)

		publishedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("error when scrape feed on published time  %v", err)
		}

		err = state.dbQueriesData.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedTime,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			return fmt.Errorf("error when scrape feed on create post index %d: %v", idx, err)
		}
	}

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

func handlerUnFollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("url argument is required")
	}

	feed, err := state.dbQueriesData.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error on handler unfollow get feed by url: %v", err)
	}

	err = state.dbQueriesData.DeleteFollow(context.Background(), database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error on handler unfollow delete follow: %v", err)
	}

	return nil
}

func handlerBrowse(state *state, cmd command) error {
	var limit int
	if len(cmd.args) < 1 {
		limit = 2
	} else {
		converted, err := strconv.Atoi(cmd.args[0])
		limit = converted
		if err != nil {
			return fmt.Errorf("error on handler browse on convert: %v", err)
		}
	}

	posts, err := state.dbQueriesData.GetPosts(context.Background(), int32(limit))
	if err != nil {
		return fmt.Errorf("error on handler browse on get post: %v", err)
	}

	for _, item := range posts {
		fmt.Printf("The title of the post %s\n", item.Title)
		fmt.Printf("Published at %s\n", item.PublishedAt)
	}

	return nil
}
