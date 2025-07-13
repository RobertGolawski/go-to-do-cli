package shared

import (
	"encoding/json"
	"os"

	"github.com/RobertGolawski/go-to-do-cli/models"
	"github.com/spf13/viper"
)

func Sync(list *models.TodoList) error {

	newContent, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		// log.Printf("Error marshalling list after mark: %v", err)
		return err
	}
	err = os.WriteFile(viper.GetString("todopath"), newContent, 0644)
	if err != nil {
		// log.Printf("Error writing to file after mark: %v", err)
		return err
	}

	return nil
}
