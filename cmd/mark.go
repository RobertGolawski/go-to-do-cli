/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/RobertGolawski/go-to-do-cli/shared"
	"github.com/spf13/cobra"
)

// markCmd represents the mark command
var markCmd = &cobra.Command{
	Use:   "mark",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mark called with args: %v\n", args[0])

		idString := args[0]

		parsedID, err := strconv.Atoi(idString)
		if err != nil {
			log.Printf("Incorrect ID, expected an int got: %v", idString)
			return
		}

		err = list.MarkDone(parsedID)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		err = shared.Sync(list)
		if err != nil {
			log.Printf("Error with sync in mark: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(markCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// markCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// markCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
