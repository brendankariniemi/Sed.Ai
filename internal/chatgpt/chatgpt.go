package chatgpt

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"os/exec"
	"sai/internal/messages"
	"strings"
)

const (
	invalidInputResponse    = "Invalid input. Please provide a valid command request or be more specific!"
	invalidCommandGenerated = "ChatGPT was unable to generate a valid command. Please try different wording!"
)

// GetResponse generates a response using ChatGPT for the given message chain.
func GetResponse(messageChain messages.MessageChain) (openai.ChatCompletionMessage, error) {
	apiKey := viper.GetString("openai_api_key")
	temperature := viper.GetFloat64("openai_temperature")
	topP := viper.GetFloat64("openai_top_p")

	client := openai.NewClient(apiKey)

	request := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    messageChain.Messages,
		Temperature: float32(temperature),
		TopP:        float32(topP),
	}

	resp, err := client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return openai.ChatCompletionMessage{}, errors.New("ChatGPT request error: " + err.Error())
	}

	if strings.Contains(resp.Choices[0].Message.Content, "InvalidInput") {
		return openai.ChatCompletionMessage{}, errors.New(invalidInputResponse)
	}

	for _, choice := range resp.Choices {
		cmd := choice.Message.Content
		words := strings.Fields(cmd)

		_, err := exec.LookPath(words[0])
		if err == nil {
			return choice.Message, nil
		}
	}

	return openai.ChatCompletionMessage{}, errors.New(invalidCommandGenerated)
}
