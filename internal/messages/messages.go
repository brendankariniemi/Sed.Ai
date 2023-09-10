package messages

import "github.com/sashabaranov/go-openai"

type MessageChain struct {
	Messages []openai.ChatCompletionMessage `json:"Messages"`
	Command  string
}

var MyMessageChain = MessageChain{
	Messages: []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are an AI-powered command assistant. Your primary function is to convert plain English input into actionable command line commands. NEVER respond with anything other than the actionable command or the text below.",
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "If you do not have enough details to form a valid command or the request is not related to generating actionable command line commands return: InvalidInput",
		},
	},
}

// SetFinalCommand sets the final command in MyMessageChain.
func SetFinalCommand(command string) {
	MyMessageChain.Command = command
}

// SetMessages sets the messages in MyMessageChain.
func SetMessages(messages []openai.ChatCompletionMessage) {
	MyMessageChain.Messages = messages
}

// AppendMessage appends a new message to MyMessageChain.
func AppendMessage(message string, role string) {
	MyMessageChain.Messages = append(MyMessageChain.Messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: message,
	})
}

// FUNCTIONS BELOW ACT ON HISTORICAL MESSAGE CHAINS

// GetFirstUserMessage returns the content of the first user message in a given message chain.
func GetFirstUserMessage(messageChain MessageChain) string {
	for _, message := range messageChain.Messages {
		if message.Role == openai.ChatMessageRoleUser {
			return message.Content
		}
	}
	return ""
}

// GetLastUserInputAndTruncate returns the content of the last user message in a message chain
// and truncates the chain to include only messages up to the last user message.
func GetLastUserInputAndTruncate(messageChain *MessageChain) string {
	var lastUserMessage string
	var lastIndex int

	for i, message := range messageChain.Messages {
		if message.Role == openai.ChatMessageRoleUser {
			lastUserMessage = message.Content
			lastIndex = i
		}
	}

	if lastIndex >= 0 && lastIndex < len(messageChain.Messages)-1 {
		messageChain.Messages = messageChain.Messages[:lastIndex]
	}

	return lastUserMessage
}
