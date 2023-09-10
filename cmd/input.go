package cmd

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
	"sai/internal/history"
	"sai/internal/messages"
)

// getNewOrHistoricalCommand gets user input and manages historical command navigation.
func getNewOrHistoricalCommand(user string) (string, int) {
	fmt.Print("Sed.Ai: What do you want to do?\n")
	fmt.Print(user + ": ")

	inputCh := make(chan string)
	quitCh := make(chan struct{})

	err := keyboard.Open()
	checkError("Error opening keyboard", err)
	defer keyboard.Close()

	historyPosition := -1
	var promptToDisplay string
	input := ""

	go func() {
		defer close(inputCh)

		for {
			select {
			case <-quitCh:
				return
			default:
				char, key, err := keyboard.GetKey()
				if err != nil {
					continue
				}

				switch key {
				case keyboard.KeyEnter:
					inputCh <- promptToDisplay
				case keyboard.KeyArrowUp:
					historyPosition++
				case keyboard.KeyArrowDown:
					historyPosition--
				case keyboard.KeyCtrlC:
					os.Exit(0)
				case keyboard.KeySpace:
					input += " "
				case keyboard.KeyBackspace, keyboard.KeyBackspace2:
					if len(input) > 0 {
						input = input[:len(input)-1]
					}
				default:
					input += string(char)
				}

				if historyPosition <= -1 || len(history.MessageChainsHistorical) == 0 {
					historyPosition = -1
					promptToDisplay = input
				} else if historyPosition >= len(history.MessageChainsHistorical) {
					historyPosition = len(history.MessageChainsHistorical) - 1
					promptToDisplay = fmt.Sprintf("{ %s ... } -> %s",
						messages.GetFirstUserMessage(history.MessageChainsHistorical[historyPosition]),
						history.MessageChainsHistorical[historyPosition].Command)
				} else {
					promptToDisplay = fmt.Sprintf("{ %s ... } -> %s",
						messages.GetFirstUserMessage(history.MessageChainsHistorical[historyPosition]),
						history.MessageChainsHistorical[historyPosition].Command)
				}

				fmt.Printf("li\r\033[K%s: %s", user, promptToDisplay)
			}
		}
	}()

	return <-inputCh, historyPosition
}

// setInitialMessageChain sets the initial message chain from historical data.
func setInitialMessageChain(historyPosition int) string {
	historicalMessages := &history.MessageChainsHistorical[historyPosition]
	userInput := messages.GetLastUserInputAndTruncate(historicalMessages)
	messages.SetMessages(historicalMessages.Messages)
	history.RemoveOldPromptsFromHistory(historyPosition)
	return userInput
}

// getCommandDetails gets additional command details from the user.
func getCommandDetails(user string) string {
	fmt.Print("Sed.Ai: Provide more details about the command.\n")
	fmt.Print(user + ": ")
	return readString()
}
