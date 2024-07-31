package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
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
				log.Fatalf("ERROR: %s", err)
			}

			fmt.Fprintf(os.Stderr, "Image generated, HEAD request on URL returned %d \n", resp.StatusCode)
			fmt.Fprintf(os.Stdout, "Image URL: %s\n", *generatedImage.URL)

			downloadFlag := cmd.Flags().Lookup("download")
			if downloadFlag != nil && downloadFlag.Changed {
				fmt.Println("Downloading image...")
				url := *generatedImage.URL

				response, err := http.Get(url)
				if err != nil {
					log.Fatal(err)
				}
				defer response.Body.Close()

				file, err := os.CreateTemp("/tmp", "*.jpg")
				if err != nil {
					log.Fatal("Error creating file:", err)
					return
				}
				defer file.Close()

				_, err = io.Copy(file, response.Body)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("Success!\nYour image has been downloaded and stored in your /tmp folder with the filename: ", file.Name())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().BoolP("download", "d", false, "download image to local device")
}
