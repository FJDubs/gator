package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/FJDubs/gator/internal/config"
	"github.com/FJDubs/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handler map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	fun, exists := c.handler[cmd.name]
	if !exists {
		return fmt.Errorf("command is not supported: %v", cmd.name)
	}
	return fun(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handler[name] = f
}

func main() {
	var s state
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v\n", err)
		os.Exit(1)
	}
	s.cfg = &cfg

	db, err := sql.Open("postgres", s.cfg.DbUrl)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.db = dbQueries

	var cmds commands
	cmds.handler = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough args provided to run a command")
		os.Exit(1)
	}
	cmd := command{args[1], args[2:]}
	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
		os.Exit(1)
	}
}
