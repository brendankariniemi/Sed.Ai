package history

import (
	"encoding/json"
	"github.com/spf13/viper"
	"os"
	"sai/internal/messages"
)

const historyFile = "/internal/history/history.json"

var MessageChainsHistorical []messages.MessageChain

// LoadHistory reads historical message chains from a file into memory.
func LoadHistory() error {
	exeDirectory := viper.GetString("exe_directory")
	file, err := openFileForReading(exeDirectory + historyFile)
	if err != nil {
		return err
	}
	defer closeFile(file)

	err = json.NewDecoder(file).Decode(&MessageChainsHistorical)
	if err != nil {
		return err
	}

	return nil
}

// UpdateAndSaveHistory appends a new message chain to the historical data and saves it to a file.
func UpdateAndSaveHistory(chain messages.MessageChain) error {
	maxHistorySize := viper.GetInt("history_max_size")

	tempMessageChains := []messages.MessageChain{chain}
	tempMessageChains = append(tempMessageChains, MessageChainsHistorical...)
	if len(tempMessageChains) > maxHistorySize {
		tempMessageChains = tempMessageChains[:maxHistorySize]
	}

	MessageChainsHistorical = tempMessageChains

	updatedData, err := json.MarshalIndent(MessageChainsHistorical, "", "  ")
	if err != nil {
		return err
	}

	exeDirectory := viper.GetString("exe_directory")
	file, err := openFileForWriting(exeDirectory + historyFile)
	if err != nil {
		return err
	}
	defer closeFile(file)

	_, err = file.Write(updatedData)
	if err != nil {
		return err
	}

	return nil
}

func RemoveOldPromptsFromHistory(historyPosition int) {
	if historyPosition < 0 || historyPosition >= len(MessageChainsHistorical) {
		return
	}

	MessageChainsHistorical = append(MessageChainsHistorical[:historyPosition], MessageChainsHistorical[historyPosition+1:]...)
}

// openFileForReading opens a file in read mode and returns the file handle.
func openFileForReading(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// openFileForWriting creates or opens a file for writing and returns the file handle.
func openFileForWriting(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// closeFile safely closes a file handle.
func closeFile(file *os.File) {
	if file != nil {
		file.Close()
	}
}
