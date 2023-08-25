package chatgpt

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"os/exec"
	"strings"
)

const (
	invalidInputResponse    = "Invalid input. Please provide a valid command request or be more specific!"
	invalidCommandGenerated = "ChatGPT was unable to generate a valid command. Please try different wording!"
	defaultTemperature      = 0.1
	defaultTopP             = 0.1
)

func GetResponse(messages []openai.ChatCompletionMessage) (openai.ChatCompletionMessage, error) {
	apiKey := viper.GetString("api_key")
	client := openai.NewClient(apiKey)

	request := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    messages,
		Temperature: defaultTemperature,
		TopP:        defaultTopP,
	}

	// Make the ChatGPT request
	resp, err := client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return openai.ChatCompletionMessage{}, errors.New("ChatGPT request error: " + err.Error())
	}

	// Check if the response contains an "InvalidInput" message
	if strings.Contains(resp.Choices[0].Message.Content, "InvalidInput") {
		return openai.ChatCompletionMessage{}, errors.New(invalidInputResponse)
	}

	// Loop through response choices to find a valid command
	for _, choice := range resp.Choices {
		cmd := choice.Message.Content
		words := strings.Fields(cmd)

		// Check if the command's executable exists in the system
		_, err := exec.LookPath(words[0])
		if err == nil {
			return choice.Message, nil
		}
	}

	// If no valid command found return: invalidCommandGenerated
	return openai.ChatCompletionMessage{}, errors.New(invalidCommandGenerated)
}
