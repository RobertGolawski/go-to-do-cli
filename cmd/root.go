/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/RobertGolawski/go-to-do-cli/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var list *models.TodoList

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-to-do-cli",
	Short: "A CLI to-do app",
	Long:  `A test app to get used to cobra, viper, bubbletea and lipgloss used to learn how to build cli apps in Go. Basic functionality is adding a todo, deleting a todo, updating a todo, completing a todo, and view todo.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		todoFilePath := viper.GetString("todopath")

		if todoFilePath == "" {
			log.Println("todo file not found, initialising empty list")
			list = &models.TodoList{Todos: []models.TodoItem{}, NextID: 0}
			return nil
		}

		if _, err := os.Stat(todoFilePath); os.IsNotExist(err) {
			log.Printf("File at %v does not exist, initialising empty list", todoFilePath)
			list = &models.TodoList{Todos: []models.TodoItem{}, NextID: 0}
			return nil
		}

		data, err := os.ReadFile(todoFilePath)
		if err != nil {
			log.Printf("Error when reading file at %v", todoFilePath)
			return fmt.Errorf("%v", err)
		}

		log.Printf("this is inside of data: %v ", string(data))

		if len(data) != 0 {
			if err := json.Unmarshal(data, &list); err != nil {
				return fmt.Errorf("Failed unmarshalling the contents off the file with error: %v", err)
			}
		} else {
			log.Printf("File was found at %v but contents was empty, initialising new list", todoFilePath)
			list = &models.TodoList{Todos: []models.TodoItem{}, NextID: 0}
		}

		log.Printf("List on load: %v", list)

		return nil
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Println("This ran post lol")
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-to-do-cli.yaml)")

	cobra.OnInitialize(initConfig)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error when trying to find user home dir in root: %v", err)
	}

	todoDir := filepath.Join(home, ".todos")
	if _, err := os.Stat(todoDir); err != nil {
		if os.IsNotExist(err) {
			mkErr := os.MkdirAll(todoDir, 0755)
			if mkErr != nil {
				log.Fatalf("Error making directory in initConfig: %v ", mkErr)
			}
		} else {
			log.Fatalf("Error getting stats on directory in initConfig: %v", err)
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home + "/.todos")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetDefault("todopath", "")
			viper.SafeWriteConfig()
		} else {
			log.Fatalf("Error reading config: %v", err)
		}
	}
	viper.AutomaticEnv()
}
