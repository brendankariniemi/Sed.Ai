package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// checkError checks for and handles errors.
func checkError(errMessage string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", errMessage, err)
		os.Exit(1)
	}
}

// readString reads and trims a string from stdin.
func readString() string {
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	checkError("Error reading input", err)
	userInput = strings.TrimSpace(userInput)
	return userInput
}
