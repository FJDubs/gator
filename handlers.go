package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/FJDubs/gator/internal/database"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("user does not exists: %v\n", err)
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		log.Printf("error setting username: %v\n", err)
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the user's name")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		fmt.Println("user already exists, please choose differenet name")
		os.Exit(1)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.args[0],
	}
	usr, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	s.cfg.SetUser(usr.Name)
	fmt.Println("user has been created")
	fmt.Printf("User ID: %v\nCreated At: %v\nUpdated At: %v\nName: %v", usr.ID, usr.CreatedAt.Time, usr.UpdatedAt, usr.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("issue reseting database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("users database has been reset")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error getting users: %v\n", err)
		os.Exit(1)
	}
	currUsr := s.cfg.CurrentUserName
	for _, usr := range users {
		if currUsr == usr.Name {
			fmt.Printf("* %v (current)\n", usr.Name)
		} else {
			fmt.Printf("* %v\n", usr.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error retrieving feed: %v", err)
	}
	b, _ := json.MarshalIndent(feed, "", "  ")
	fmt.Println(string(b))
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("the addfeed handler expects a two argument, feed title and url")
	}
	usr, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    uuid.NullUUID{UUID: usr.ID, Valid: true},
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}
	s.cfg.SetUser(feed.Name)
	fmt.Println("feed has been created")
	fmt.Printf("Feed ID: %v\nCreated At: %v\nUpdated At: %v\nName: %v\nURL: %v\nUser ID: %v\n", feed.ID, feed.CreatedAt.Time, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from database: %v", err)
	}

	for _, feed := range feeds {
		usr_id := feed.UserID.UUID
		usr_name, err := s.db.GetUserNameFromUUID(context.Background(), usr_id)
		if err != nil {
			return fmt.Errorf("error retrieving user name from user id: %v", err)
		}
		fmt.Printf("Feed Name: %v\nFeed URL: %v\nFeed's Creator: %v\n\n", feed.Name, feed.Url, usr_name)
	}

	return nil
}
