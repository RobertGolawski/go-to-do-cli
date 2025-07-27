/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/RobertGolawski/go-to-do-cli/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")

		piIp := viper.GetString("piip")

		if piIp == "" {
			log.Println("no api address found")
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%vpull", piIp), nil)
		if err != nil {
			log.Printf("error making the get request: %v", err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error during the request: %v", err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("expected ok but got: %v", resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error reading resp body: %v", err)
			return
		}

		var list models.TodoList
		err = json.Unmarshal(body, &list)
		if err != nil {
			log.Printf("error unmarshalling file: %v", err)
			return
		}

		log.Printf("successfully pulled, saving data")

		jsonData, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			log.Printf("error marshalling data: %v", err)
			return
		}

		filePath := viper.GetString("todopath")

		if filePath == "" {
			log.Println("local filepath not set")
			return
		}

		err = os.WriteFile(filePath, jsonData, 0644)
		if err != nil {
			log.Printf("error writing to file: %v", err)
			return
		}

		log.Printf("pull successful")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
