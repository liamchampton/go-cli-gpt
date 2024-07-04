package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// exampleCmd represents the example command
var questionCmd = &cobra.Command{
	Use:   "question",
	Short: "ask gpt a question",
	Long:  `use this command to ask gpt a question and get an answer in your terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		// cmd.Help()

		// Load the .env file
		if err := godotenv.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading .env file\n")
			return
		}

		// Set environment variables
		azureOpenAIKey := os.Getenv("AOAI_COMPLETIONS_API_KEY")
		modelDeployment := os.Getenv("AOAI_COMPLETIONS_MODEL")
		azureOpenAIEndpoint := os.Getenv("AOAI_COMPLETIONS_ENDPOINT")

		if azureOpenAIKey == "" || modelDeployment == "" || azureOpenAIEndpoint == "" {
			fmt.Fprintf(os.Stderr, "Skipping command, environment variables missing\n")
			return
		}

		// Get question from user input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your question: ")
		question, _ := reader.ReadString('\n')
		fmt.Println("You entered:", question)

		keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

		client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

		if err != nil {
			// TODO: Update the following line with specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		resp, err := client.GetCompletions(context.TODO(), azopenai.CompletionsOptions{
			// Prompt: []string{"What is Azure OpenAI, in 20 words or less"},
			Prompt:         []string{question},
			MaxTokens:      to.Ptr(int32(2048)),
			Temperature:    to.Ptr(float32(0.0)),
			DeploymentName: &modelDeployment,
		}, nil)

		if err != nil {
			// TODO: Update the following line with specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		for _, choice := range resp.Choices {
			fmt.Fprintf(os.Stderr, "Result: %s\n", *choice.Text)
		}

	},
}

func init() {
	rootCmd.AddCommand(questionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
