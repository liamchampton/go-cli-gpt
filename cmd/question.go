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

		azureOpenAIKey := os.Getenv("AZURE_OPENAI_API_KEY")
		modelDeploymentID := os.Getenv("YOUR_MODEL_DEPLOYMENT_NAME")
		maxTokens := int32(400)
		azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")

		if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
			fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
			return
		}

		keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

		// Get question from user input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your question: ")
		question, _ := reader.ReadString('\n')
		fmt.Println("You entered:", question)

		client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

		if err != nil {
			// TODO: Update with application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		// NOTE: all messages, regardless of role, count against token usage for this API.
		messages := []azopenai.ChatRequestMessageClassification{
			// You set the tone and rules of the conversation with a prompt as the system role.
			&azopenai.ChatRequestSystemMessage{Content: to.Ptr("You are a helpful assistant to help with all sorts of user questions. If you do not know the answer, then your response must be 'I am not smart enough to answer that question. Please try again or ask ChatGPT instead.'.")},

			// The user asks a question
			// &azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("Does Azure OpenAI support customer managed keys?")},
			&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(question)},

			// The reply would come back from the model. You'd add it to the conversation so we can maintain context.
			// &azopenai.ChatRequestAssistantMessage{Content: to.Ptr("Yes, customer managed keys are supported by Azure OpenAI")},

			// The user answers the question based on the latest reply.
			// &azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("What other Azure Services support customer managed keys?")},

			// from here you'd keep iterating, sending responses back from ChatGPT
		}

		gotReply := false

		resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
			// NOTE: all messages count against token usage for this API.
			Messages:       messages,
			DeploymentName: &modelDeploymentID,
			MaxTokens:      &maxTokens,
		}, nil)

		if err != nil {
			// TODO: Update with application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		for _, choice := range resp.Choices {
			gotReply = true

			if choice.ContentFilterResults != nil {
				fmt.Fprintf(os.Stderr, "Content filter results\n")

				if choice.ContentFilterResults.Error != nil {
					fmt.Fprintf(os.Stderr, "  Error:%v\n", choice.ContentFilterResults.Error)
				}

				fmt.Fprintf(os.Stderr, "  Hate: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Hate.Severity, *choice.ContentFilterResults.Hate.Filtered)
				fmt.Fprintf(os.Stderr, "  SelfHarm: sev: %v, filtered: %v\n", *choice.ContentFilterResults.SelfHarm.Severity, *choice.ContentFilterResults.SelfHarm.Filtered)
				fmt.Fprintf(os.Stderr, "  Sexual: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Sexual.Severity, *choice.ContentFilterResults.Sexual.Filtered)
				fmt.Fprintf(os.Stderr, "  Violence: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Violence.Severity, *choice.ContentFilterResults.Violence.Filtered)
			}

			if choice.Message != nil && choice.Message.Content != nil {
				fmt.Fprintf(os.Stderr, "Content[%d]: %s\n", *choice.Index, *choice.Message.Content)
			}

			if choice.FinishReason != nil {
				// this choice's conversation is complete.
				fmt.Fprintf(os.Stderr, "Finish reason[%d]: %s\n", *choice.Index, *choice.FinishReason)
			}
		}

		if gotReply {
			fmt.Fprintf(os.Stderr, "Received chat completions reply\n")
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
