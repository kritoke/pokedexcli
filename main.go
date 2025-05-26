package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	userInput := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		userInput.Scan()
		firstWord := cleanInput(userInput.Text())
		fmt.Println("Your command was:", firstWord[0])
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
