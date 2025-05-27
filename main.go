package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	var commandNames []string
	for name := range commands {
		commandNames = append(commandNames, name)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(commandNames)))

	for _, cmd := range commandNames {
		cmd := commands[cmd]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func main() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}

	userInput := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		userInput.Scan()
		input := cleanInput(userInput.Text())

		if len(input) == 0 || input[0] == "" {
			continue
		}

		command := input[0]

		if cmd, exists := commands[command]; exists {
			err := cmd.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(strings.TrimSpace(text))

	if text == "" {
		return []string{""}
	}

	fields := strings.Fields(text)

	if len(fields) == 0 {
		return []string{""}
	}

	return fields
}
