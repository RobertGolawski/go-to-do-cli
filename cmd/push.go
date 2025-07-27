/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push to save remotely",
	Long:  `Push to a remote api to save on the api host system`,
	Run: func(cmd *cobra.Command, args []string) {
		piIP := viper.GetString("piip")
		todopath := viper.GetString("todopath")

		if piIP == "" || todopath == "" {
			log.Printf("one of ip :%v or todopath: %v is empty", piIP, todopath)
			return
		}

		if _, err := os.Stat(todopath); err != nil {
			if os.IsNotExist(err) {
				log.Printf("no file found at %v", todopath)
				return
			}
		}

		file, err := os.Open(todopath)
		if err != nil {
			log.Printf("error opening file: %v", err)
			return
		}

		defer file.Close()

		req, err := http.NewRequest("POST", fmt.Sprintf("%vpush", piIP), file)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error sending request: %v", err)
			return
		}

		log.Printf("success, resp: %v", resp)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
