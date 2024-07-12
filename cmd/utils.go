package cmd

import (
	"bufio"
	"fmt"
	"os"
)

func GetUserInput(userPrint string) string {
	reader := bufio.NewReader(os.Stdin)
	// fmt.Print("I am a generic helper. I may or may not know the answer to your question so be patient with me. \nPlease enter your question: ")
	fmt.Print(userPrint)
	userInput, _ := reader.ReadString('\n')
	return userInput
}
