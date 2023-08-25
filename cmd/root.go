package cmd

import (
	"bufio"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"sai/internal/chatgpt"
	"sai/internal/config"
	"strings"
)

var messages = []openai.ChatCompletionMessage{
	{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are an AI-powered command assistant. Your primary function is to convert plain English input into actionable command line commands. NEVER respond with anything other than the actionable command or the text below.",
	},
	{
		Role:    openai.ChatMessageRoleSystem,
		Content: "If you do not have enough details to form a valid command or the request is not related to generating actionable command line commands return: InvalidInput",
	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Sed.Ai [text]",
	Short: "Sed.Ai is a CLI tool that leverages ChatGPT to convert plain text input into executable command line commands.",
	Run: func(cmd *cobra.Command, args []string) {
		var userInput string
		user := viper.GetString("username")

		// Check if any arguments were provided
		if len(args) == 0 {
			// No arguments provided, prompt for input
			fmt.Print("Sed.Ai: What do you want to do?\n")
			fmt.Print(user + ": ")
			userInput = readString()
		} else {
			userInput = args[0]
		}

		for {
			// Add user input to message history
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userInput,
			})

			// Call chatgpt
			message, err := chatgpt.GetResponse(messages)
			if err != nil {
				cmd.PrintErrf("Error generating command: %v\n", err)
				os.Exit(1)
			}

			// Add chatgpt response to chat history
			messages = append(messages, message)

			// Display response, and check if we should execute it
			cmd.Printf("Sed.Ai: %s\n", message.Content)
			cmd.Print("Sed.Ai: Would you like to execute this command? (Y/n): ")
			userInput = readString()

			// Either execute or get more details
			if userInput == "Y" || userInput == "y" || userInput == "" {
				cmd := exec.Command(message.Content)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err = cmd.Run()
				if err != nil {
					// Handle command execution error
					fmt.Printf("Error executing command: %v\n", err)
					os.Exit(2)
				}
				break
			} else {
				fmt.Print("Sed.Ai: Provide more details about the command.\n")
				fmt.Print(user + ": ")
				userInput = readString()
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once for the rootCmd.
func Execute() {
	config.InitConfig()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Sed.Ai.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func readString() string {
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}
	userInput = strings.TrimSpace(userInput)
	return userInput
}
