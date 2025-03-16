package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/tuananhlai/brevity-go/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Use API key from config, fallback to environment variable
	apiKey := cfg.DeepseekAPIKey
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY not set in config or environment")
	}

	client := openai.NewClient(
		option.WithBaseURL("https://api.deepseek.com"),
		option.WithAPIKey(apiKey),
	)

	// Create a new chat completion
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a front-desk helper at a 5-star hotel."),
			openai.UserMessage("Can you tell me the way to my room? I'm in room 818."),
		}),
		Model:               openai.F("deepseek-chat"),
		MaxCompletionTokens: openai.F(int64(200)),
	})
	if err != nil {
		log.Fatalf("Error creating chat completion: %v", err)
	}

	fmt.Println("totalTokens", chatCompletion.Usage.TotalTokens)
	fmt.Println(chatCompletion.Choices[0].Message.Content)
}
