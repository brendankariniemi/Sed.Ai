package cmd

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sai/internal/chatgpt"
	"sai/internal/config"
	"sai/internal/history"
	"sai/internal/messages"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Sed.Ai [text]",
	Short: "Sed.Ai is a CLI tool that leverages ChatGPT to convert plain text input into executable command line commands.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the executing user
		user := viper.GetString("username")

		// Load previously executed commands
		err := history.LoadHistory()
		checkError("Error loading history", err)

		// Set the initial message chain
		userInput, historyIndex := getNewOrHistoricalCommand(user)
		if historyIndex >= 0 {
			userInput = setInitialMessageChain(historyIndex)
		}
		fmt.Print("\n")

		// Main processing loop
		for {
			// Add user input to message history
			messages.AppendMessage(userInput, openai.ChatMessageRoleUser)

			// Call chatgpt
			message, err := chatgpt.GetResponse(messages.MyMessageChain)
			checkError("Error generating command", err)

			// Add chatgpt response to message history
			messages.AppendMessage(message.Content, openai.ChatMessageRoleSystem)

			// Display response, and check if we should execute it
			if shouldExecuteCommand(message.Content) {
				executeCommand(message.Content)

				// Save the command in history
				err = history.UpdateAndSaveHistory(messages.MyMessageChain)
				checkError("Error saving history", err)
				break
			} else {
				// Prompt user for more details about command
				userInput = getCommandDetails(user)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once for the rootCmd.
func Execute() {
	config.InitConfig()
	err := rootCmd.Execute()
	checkError("Error running root command", err)
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
