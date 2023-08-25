package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
)

func InitConfig() {
	// Get the path to the symlink
	symPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	// Resolve the symbolic link to get the path of the actual executable
	exePath, err := filepath.EvalSymlinks(symPath)
	if err != nil {
		fmt.Printf("Error resolving symbolic link: %v\n", err)
		os.Exit(1)
	}

	// Set the configuration file name and path relative to the executable
	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Dir(exePath) + "/internal/config")

	// Attempt to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	// Retrieve the current user's information
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
		return
	}

	// Store the username in Viper for later use
	viper.Set("username", currentUser.Username)
}
