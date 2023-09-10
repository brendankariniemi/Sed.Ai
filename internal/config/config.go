package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
)

func InitConfig() {
	symPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	exePath, err := filepath.EvalSymlinks(symPath)
	if err != nil {
		fmt.Printf("Error resolving symbolic link: %v\n", err)
		os.Exit(1)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Dir(exePath) + "/internal/config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
		return
	}

	viper.Set("exe_directory", filepath.Dir(exePath))
	viper.Set("username", currentUser.Username)
}
