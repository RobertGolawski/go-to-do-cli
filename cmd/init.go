/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var filePathFlag string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises the local storage file for todos",
	Long:  `Initialises the local storage by making the storage at ~/.todos/go-cli-todos.json creating any necessary parent directories. Only works if the file doesn't already exist.`,
	Run: func(cmd *cobra.Command, args []string) {

		if 1 < len(args) {
			log.Fatalf("Too many arguments, len: %v", len(args))
		}

		existingPath := viper.GetString("todopath")

		if existingPath != "" {
			if _, err := os.Stat(existingPath); err == nil {
				log.Printf("A file with ToDos already exists at: %v . Please remove it first before initialising a second one.", existingPath)
				return
			}
		}

		fmt.Println("init called")
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Encountered an error during checking home directory: %v", err)
		}

		var filePath string
		if filePathFlag != "" {
			if err := validatePath(home, filePathFlag); err != nil {
				log.Fatalf("Invalid directory path: %v", err)
			}
			filePath = filepath.Join(filePathFlag, "cli-todos.json")

		} else {
			filePath = filepath.Join(home, ".todos", "cli-todos.json")
		}
		_, err = os.Stat(filePath)
		if err == nil {
			log.Printf("File already exists at %v", filePath)
			return
		}

		dir := filepath.Dir(filePath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Error making directories: %v", err)
		}

		f, fErr := os.Create(filePath)
		if fErr != nil {
			log.Fatalf("Error creating file: %v", fErr)
		}
		defer f.Close()

		viper.Set("todopath", filePath)
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Error writing to config file: %v", err)
		}

		log.Printf("File successfully created at: %v", filePath)
	},
}

func validatePath(home, path string) error {

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("Error when getting absolute target path: %v", err)
	}

	forbidden := []string{"\x00", "?", "%", "*", ":", "|", "\"", "<", ">"}
	for _, c := range forbidden {
		if strings.Contains(absPath, c) {
			return errors.New("Path contains forbidden characters")
		}
	}

	absHome, hErr := filepath.Abs(home)
	if hErr != nil {
		return fmt.Errorf("Error when getting absolute home path: %v", hErr)
	}

	rel, err := filepath.Rel(absHome, absPath)
	if err != nil {
		return errors.New("Failed to compute relative path")
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return errors.New("Path is outside of permitted directory")
	}

	return nil

}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&filePathFlag, "dir", "d", "", "Custom directory for the todos JSON file")

}
