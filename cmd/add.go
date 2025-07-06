/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/RobertGolawski/go-to-do-cli/models"
	"github.com/spf13/cobra"
)

var dueDate string
var toDo string
var prio string
var daysOfTheWeek = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("add called with args: %v", args)

		task := args[0]

		var d time.Time
		var err error
		if dueDate == "" {
			d = time.Now().Truncate(24 * time.Hour)
		} else {
			d, err = parseDate(dueDate)
			if err != nil {
				log.Printf("Error parsing date: %v", err)
				return
			}
			log.Printf("The date returned: %v", d)
		}

		toDo := models.TodoItem{
			ID:       list.NextID,
			ToDo:     task,
			DueDate:  d,
			Priority: models.Medium,
			Done:     false,
		}

		log.Printf("This todo is ready to be added: %v", toDo)

	},
}

func parseDate(input string) (time.Time, error) {
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "today":
		return time.Now().Truncate(24 * time.Hour), nil
	case "tomorrow":
		return time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour), nil
	}

	if isWeekday(input) {
		return parseWeekday(input), nil
	}

	ret, err := parseNumericDate(input)
	if err != nil {
		log.Printf("Failed parsing numeric date: %v", input)
		return time.Time{}, err
	}

	return ret, nil
}

func isWeekday(input string) bool {
	_, exists := daysOfTheWeek[input]
	return exists
}

func parseWeekday(input string) time.Time {
	wd, ok := daysOfTheWeek[input]
	if !ok {
		log.Fatalf("Error getting weekday from map")
	}

	daysAhead := int(wd - time.Now().Weekday())
	if daysAhead <= 0 {
		daysAhead += 7
	}
	return time.Now().AddDate(0, 0, daysAhead)
}

func parseNumericDate(input string) (time.Time, error) {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == '/' || r == '-' || r == '.' || r == ' '
	})

	switch len(parts) {
	case 1:
		return parseDay(parts[0])
	case 2:
		return parseDayMonth(parts[0], parts[1])
	case 3:
		return parseFullDate(parts[0], parts[1], parts[2])
	default:
		return time.Time{}, fmt.Errorf("Invalid date format: %v", input)
	}
}

func parseDay(day string) (time.Time, error) {
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		log.Fatalf("Invalid day argument: %v caused error: %v", day, err)
	}
	dayInMonth := time.Date(time.Now().Year(), time.Now().Month(), dayInt, 0, 0, 0, 0, time.Local)
	dayNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	if dayInMonth.Before(dayNow) {
		log.Printf("The day %v has already passed this month, adding for next month", day)
		return dayInMonth.AddDate(0, 1, 0), nil
	} else if dayInMonth.After(dayNow) {
		return dayInMonth, nil
	} else {
		return parseDate("today")
	}
}

func parseDayMonth(day, month string) (time.Time, error) {

	var dayInt int
	var monthInt int
	var err error

	if dayInt, err = strconv.Atoi(day); err != nil {
		log.Fatalf("Invalid day argument: %v caused error: %v", day, err)
	}
	if monthInt, err = strconv.Atoi(month); err != nil {
		log.Fatalf("Invalid month argument: %v caused error: %v", month, err)
	}

	dayMonthInYear := time.Date(time.Now().Year(), time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)
	dayMonthNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	if dayMonthInYear.Before(dayMonthNow) {
		dayInMonth := time.Date(time.Now().Year(), time.Now().Month(), dayInt, 0, 0, 0, 0, time.Local)
		if dayInMonth.Before(dayMonthNow) {
			log.Printf("Provided %v as a month and %v as the day, the month has already passed and the day has already passed in current month. Adding todo for %v in the next month", month, day, day)
			return dayInMonth.AddDate(0, 1, 0), nil
		}
		log.Printf("The month %v is in the past, adding the todo for the %v of the current month", month, day)
		return time.Date(time.Now().Year(), time.Now().Month(), dayInt, 0, 0, 0, 0, time.Local), nil
	} else if dayMonthInYear.After(dayMonthNow) {
		return dayMonthInYear, nil
	} else {
		return parseDate("today")
	}
}

func parseFullDate(day, month, year string) (time.Time, error) {

	var dayInt int
	var monthInt int
	var yearInt int
	var err error

	if dayInt, err = strconv.Atoi(day); err != nil {
		log.Fatalf("Invalid day argument: %v caused error: %v", day, err)
	}
	if monthInt, err = strconv.Atoi(month); err != nil {
		log.Fatalf("Invalid month argument: %v caused error: %v", month, err)
	}
	if yearInt, err = strconv.Atoi(year); err != nil {
		log.Fatalf("Invalid year argument: %v caused error: %v", year, err)
	}
	fullDate := time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	if fullDate.Before(now) {
		return time.Date(time.Now().Year(), time.Now().Month(), dayInt, 0, 0, 0, 0, time.Local), nil
	} else if fullDate.After(now) {
		return fullDate, nil
	} else {
		return parseDate("today")
	}
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&dueDate, "date", "d", "", "Custom due date for when the todo is meant to be completed by. Defaults to today.")
	addCmd.Flags().StringVarP(&prio, "prio", "p", "", "The priority of how important the todo is. Defaults to medium.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
