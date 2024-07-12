package cmd

import (
	"bufio"
	"fmt"
	"os"
)

func GetUserInput(userPrint string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(userPrint)
	userInput, _ := reader.ReadString('\n')
	return userInput
}
