package cmd

import (
	"context"
	"encoding/json"
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
var weatherCmd = &cobra.Command{
	Use:   "get-weather",
	Short: "Get current weather information",
	Long:  `Get current weather information for a location using function calling`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the .env file
		if err := godotenv.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading .env file\n")
			return
		}

		azureOpenAIKey := os.Getenv("AZURE_OPENAI_API_KEY")
		modelDeploymentID := os.Getenv("YOUR_MODEL_DEPLOYMENT_NAME")

		// Ex: "https://<your-azure-openai-host>.openai.azure.com"
		azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")

		if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
			fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
			return
		}

		keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

		client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
			DeploymentName: &modelDeploymentID,
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestUserMessage{
					Content: azopenai.NewChatRequestUserMessageContent("What's the weather like in London in the UK? Give this to me in celsius"),
				},
			},
			Tools: []azopenai.ChatCompletionsToolDefinitionClassification{
				&azopenai.ChatCompletionsFunctionToolDefinition{
					Function: &azopenai.FunctionDefinition{
						Name:        to.Ptr("get_current_weather"),
						Description: to.Ptr("Get the current weather in a given location"),
						Parameters: map[string]any{
							"required": []string{"location"},
							"type":     "object",
							"properties": map[string]any{
								"location": map[string]any{
									"type":        "string",
									"description": "The city and country, e.g. London, UK",
								},
								"unit": map[string]any{
									"type": "string",
									"enum": []string{"celsius", "fahrenheit"},
								},
							},
						},
					},
				},
			},
			Temperature: to.Ptr[float32](0.0),
		}, nil)

		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		funcCall := resp.Choices[0].Message.ToolCalls[0].(*azopenai.ChatCompletionsFunctionToolCall).Function

		// This is the function name we gave in the call to GetCompletions
		// Prints: Function name: "get_current_weather"
		fmt.Fprintf(os.Stderr, "Function name: %q\n", *funcCall.Name)

		// The arguments for the function come back as a JSON string
		// The arguments are pulled from the natural language query
		var funcParams *struct {
			Location string `json:"location"`
			Unit     string `json:"unit"`
		}
		err = json.Unmarshal([]byte(*funcCall.Arguments), &funcParams)

		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		// Prints:
		// Parameters: azopenai_test.location{Location:"London, UK", Unit:"celsius"}
		fmt.Fprintf(os.Stderr, "Parameters: %#v\n", *funcParams)

	},
}

func init() {
	rootCmd.AddCommand(weatherCmd)
}
