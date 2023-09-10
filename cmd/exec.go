package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sai/internal/messages"
	"strings"
)

// shouldExecuteCommand prompts the user whether they want to execute a command.
func shouldExecuteCommand(command string) bool {
	fmt.Printf("Sed.Ai: %s\n", command)
	fmt.Print("Sed.Ai: Would you like to execute this command? (Y/n): ")
	userInput := readString()
	return userInput == "Y" || userInput == "y" || userInput == ""
}

// executeCommand executes a given command.
func executeCommand(command string) {
	parts := strings.Fields(command)

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		checkError("Error executing command", err)
		return
	}

	messages.SetFinalCommand(command)
}
