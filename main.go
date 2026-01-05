package main

// The underscore tells Go that you're importing it for its side effects, not because you need to use it.
import (
	"bootDevGoRss/internal/database"
	"database/sql"

	_ "github.com/lib/pq"
)

import (
	"bootDevGoRss/internal/config"
	"fmt"
	"log"
	"os"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	configData, err := config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(configData)

	db, err := sql.Open("postgres", configData.DbUrl)
	if err != nil {
		log.Fatalf("Error connect to database %v", err)
	}
	dbQueries := database.New(db)

	stateData := state{
		configData:    &configData,
		dbQueriesData: dbQueries,
	}

	commandsData := commands{
		mapper: make(map[string]func(*state, command) error),
	}
	if err := commandsData.register("login", handlerLogin); err != nil {
		log.Fatalf("error in login command: %v", err)
	}
	if err := commandsData.register("register", handlerRegister); err != nil {
		log.Fatalf("error in register command: %v", err)
	}
	if err := commandsData.register("reset", handlerDelete); err != nil {
		log.Fatalf("error in reset command: %v", err)
	}
	if err := commandsData.register("users", handlerGetUsers); err != nil {
		log.Fatalf("error in users command: %v", err)
	}
	if err := commandsData.register("agg", handlerAggCommand); err != nil {
		log.Fatalf("error in agg command: %v", err)
	}

	// Why two? The first argument is automatically the program name, which we ignore, and we require a command name.
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	commandData := command{
		command: cmdName,
		args:    cmdArgs,
	}

	if err := commandsData.run(&stateData, commandData); err != nil {
		log.Fatal(err)
	}

	fmt.Println(configData)
}
