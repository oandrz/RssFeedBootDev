package main

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

	stateData := state{
		configData: &configData,
	}
	commandsData := commands{
		mapper: make(map[string]func(*state, command) error),
	}
	if err := commandsData.register("login", handlerLogin); err != nil {
		panic(err)
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
