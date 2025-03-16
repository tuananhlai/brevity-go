package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/tuananhlai/brevity-go/internal/article"
	"github.com/tuananhlai/brevity-go/internal/config"
)

func runGenerateArticle() {
	cfg := config.MustLoadConfig()

	apiKey := cfg.DeepseekAPIKey
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY not set in config or environment")
	}

	client := openai.NewClient(
		option.WithBaseURL("https://api.deepseek.com"),
		option.WithAPIKey(apiKey),
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a Japanese teacher who often writes various articles about " +
				"learning Japanese. Here are the non-exhaustive list of example topics you might write about: " +
				"commonly used grammar and vocabulary, vocabulary based on topics, differences between Japanese " +
				"dialects, common mistakes people make when learning Japanese, etc. Your target audience is people " +
				"learning Japanese at N1-N2 level. You write your articles in Japanese at a level that your target " +
				"audience can understand."),
			openai.UserMessage("Choose an unique and interesting topic and write an article about it."),
		}),
		Model:       openai.F("deepseek-chat"),
		Temperature: openai.F(0.5),
		MaxTokens:   openai.F(int64(200)),
	})
	if err != nil {
		log.Fatalf("Error creating chat completion: %v", err)
	}

	fmt.Println("totalTokens", chatCompletion.Usage.TotalTokens)
	fmt.Println(chatCompletion.Choices[0].Message.Content)

	// Connect to database
	db := sqlx.MustConnect("postgres", cfg.Database.URL)
	articleRepo := article.NewRepository(db)

	// Use a fixed author ID for now
	authorID := uuid.MustParse("41dc81d2-97e8-41c8-a3bf-d98322302e5c")

	// Create the article
	err = articleRepo.Create(&article.Article{
		Slug:        "my-article-slug-3921",
		Title:       "My Article Title",
		Description: "My Article Description",
		TextContent: chatCompletion.Choices[0].Message.Content,
		Content:     chatCompletion.Choices[0].Message.Content,
		AuthorID:    authorID,
	})
	if err != nil {
		log.Fatalf("Failed to create article: %v", err)
	}
}
