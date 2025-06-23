package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     *string
	Previous *string
}

var commands map[string]cliCommand

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
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

func fetchAndDisplayLocations(cfg *config, forward bool) error {
	var url string

	type Location struct {
		Count    int    `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}

	if forward {
		if cfg.Next != nil {
			url = *cfg.Next
		} else {
			url = "https://pokeapi.co/api/v2/location-area/"
		}
	} else {
		if cfg.Previous != nil {
			url = *cfg.Previous
		} else {
			fmt.Println("No previous locations, you are on the first page")
			return nil
		}
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close() // Explicitly ignore the error
	}()

	if res.StatusCode > 299 {
		return fmt.Errorf("HTTP error: %d", res.StatusCode)
	}

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	locations := Location{}
	err = json.Unmarshal(jsonData, &locations)
	if err != nil {
		return err
	}

	for _, location := range locations.Results {
		fmt.Println(string(location.Name))
	}

	cfg.Next = &locations.Next
	cfg.Previous = &locations.Previous

	return nil
}

func commandMap(cfg *config) error {
	return fetchAndDisplayLocations(cfg, true)
}

func commandMapb(cfg *config) error {
	return fetchAndDisplayLocations(cfg, false)
}

func main() {
	cfg := &config{}

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
		"map": {
			name:        "map",
			description: "Displays a map of the Pokemon locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays a map of the previous Pokemon locations",
			callback:    commandMapb,
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
			err := cmd.callback(cfg)
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
