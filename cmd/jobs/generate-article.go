package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/invopop/jsonschema"
	"github.com/jmoiron/sqlx"
	"github.com/k3a/html2text"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/tuananhlai/brevity-go/internal/config"
	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/repository"
)

type LLMOutput struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func RunGenerateArticle() {
	globalCtx := context.Background()
	cfg := config.MustLoadConfig()

	client := openai.NewClient(
		option.WithBaseURL(cfg.LLM.BaseURL),
		option.WithAPIKey(cfg.LLM.APIKey),
	)

	type Output struct {
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}

	outputSchema := generateSchema[Output]()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   "Output",
		Schema: outputSchema,
		Strict: openai.Bool(true),
	}

	chatCompletion, err := client.Chat.Completions.New(globalCtx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are a Japanese teacher who often writes various articles about learning Japanese.
				Here are the non-exhaustive list of example topics you might write about:
				commonly used grammar and vocabulary, vocabulary based on topics,
				differences between Japanese dialects, common mistakes people make when learning Japanese, etc.
				Your target audience is people learning Japanese at N1-N2 level.
				You write your articles in Japanese at a level that your target audience can understand.
				When the user asks you to write an article, you will return your answer as a JSON object, without any comment,
				adhering to the following schema.

				type Output = {
				    // The slug of the article, used for the URL. It must contains only alphanumeric characters and hyphens. Example: "my-article-slug-3921"
				    slug: string;
					// The title of the article. It should be a plain string.
					title: string;
					// A short description of the article content. Limit to 200 words or fewer.
					description: string;
					// The content of the article as valid, standard, unstyled HTML. Only these HTML tags are allowed: h2, h3, h4, p, a, img, strong, b, em, i, del,
					// strike, blockquote, pre, code, ul, ol, li, hr, br, table, thead, tbody, tr, th, td. Do not include newlines unless absolute necessary. Article should be 800 - 1500 words.
					// Keep in mind that an h1 tag with the title will be prepended afterward.
					content: string;
				}
				`),
			openai.UserMessage("Choose an unique and interesting topic and write an article about it."),
		},
		Model:       openai.ChatModel(cfg.LLM.ModelID),
		Temperature: openai.Float(2),
		MaxTokens:   openai.Int(8192),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})
	if err != nil {
		log.Fatalf("Error creating chat completion: %v", err)
	}

	fmt.Println("totalTokens", chatCompletion.Usage.TotalTokens)
	fmt.Println(chatCompletion.Choices[0].Message.Content)

	var output LLMOutput
	if err := json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &output); err != nil {
		log.Fatalf("Failed to unmarshal LLM output: %v", err)
	}

	// Connect to database
	db := sqlx.MustConnect("postgres", cfg.Database.URL)
	articleRepo := repository.NewArticleRepository(db)

	// Use a fixed author ID for now
	authorID := uuid.MustParse("2bc5d7c4-5702-4ac7-b951-c8ed0921896e")

	// Create the article
	err = articleRepo.Create(globalCtx, &model.Article{
		Slug:             output.Slug,
		Title:            output.Title,
		Description:      output.Description,
		PlaintextContent: html2text.HTML2Text(output.Content),
		Content:          output.Content,
		AuthorID:         authorID,
	})
	if err != nil {
		log.Fatalf("Failed to create article: %v", err)
	}
}

func generateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: true,
		DoNotReference:            true,
	}

	var v T
	schema := reflector.Reflect(v)

	return schema
}
