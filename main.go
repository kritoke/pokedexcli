package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Printf("Hello, World!")
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)

	if text == "" {
		return []string{""}
	}

	fields := strings.Fields(text)

	if len(fields) == 0 {
		return []string{""}
	}

	return fields
}
