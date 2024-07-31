package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func GetUserInput(userPrint string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(userPrint)
	userInput, _ := reader.ReadString('\n')
	return userInput
}

func GetLocalModel() (string, error) {
	options := []string{"llama3.1", "phi3", "mistral"}

	var selectedOption string
	prompt := &survey.Select{
		Message: "Choose a local model to use (you must have it installed!):",
		Options: options,
	}
	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting local model: %v\n", err)
		return "", err
	}

	fmt.Printf("You selected: %s\n", selectedOption)
	return selectedOption, nil
}
