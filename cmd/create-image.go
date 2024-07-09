package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var image = &cobra.Command{
	Use:   "image",
	Short: "Create an image from a prompt",
	Long:  `Create an image from a prompt using the OpenAI API and DALLE<X> model`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := godotenv.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading .env file\n")
			return
		}

		deploymentName := os.Getenv("DALLE_MODEL_NAME")
		azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
		azureOpenAIKey := os.Getenv("AZURE_OPENAI_API_KEY")

		keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

		// Get question from user input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("What image do you want to create? ")
		prompt, _ := reader.ReadString('\n')
		fmt.Println("Creating image based on your prompt... ", prompt)

		client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

		if err != nil {
			// TODO: Update with application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		resp, err := client.GetImageGenerations(context.TODO(), azopenai.ImageGenerationOptions{
			Prompt:         to.Ptr(prompt),
			ResponseFormat: to.Ptr(azopenai.ImageGenerationResponseFormatURL),
			DeploymentName: &deploymentName,
		}, nil)

		if err != nil {
			//  TODO: Update the following line with application-specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		for _, generatedImage := range resp.Data {
			// use 'azopenai.ImageGenerationResponseFormatURL'
			resp, err := http.Head(*generatedImage.URL)

			if err != nil {
				//  TODO: Update the following line with your application-specific error handling logic
				log.Fatalf("ERROR: %s", err)
			}

			fmt.Fprintf(os.Stderr, "Image generated, HEAD request on URL returned %d \n", resp.StatusCode)
			fmt.Fprintf(os.Stdout, "Image URL: %s\n", *generatedImage.URL)
		}
	},
}

func init() {
	rootCmd.AddCommand(image)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
