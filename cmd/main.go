package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")

	client := openai.NewClient(
		option.WithBaseURL("https://api.deepseek.com"),
		option.WithAPIKey(apiKey),
	)

	// Create a new chat completion
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a professional journalist who investigates financial crimes."),
			openai.UserMessage("Write an opinion piece about whether blockchain and cryptocurrency is a good financial investment, drawing on past events about these type of projects."),
		}),
		Model:               openai.F("deepseek-chat"),
		MaxCompletionTokens: openai.F(int64(400)),
	})
	if err != nil {
		log.Fatalf("Error creating chat completion: %v", err)
	}

	fmt.Println("totalTokens", chatCompletion.Usage.TotalTokens)
	fmt.Println(chatCompletion.Choices[0].Message.Content)
}
