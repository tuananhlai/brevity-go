package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/invopop/jsonschema"
	"github.com/jmoiron/sqlx"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/tuananhlai/brevity-go/internal/repository"
)

type LLMOutput struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type Article struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func RunGenerateArticle() {
	ctx := context.Background()
	baseURL := "https://openrouter.ai/api/v1"
	apiKey := ""
	client := openai.NewClient(option.WithBaseURL(baseURL), option.WithAPIKey(apiKey))

	databaseURL := "postgres://postgres:postgres@postgres:5432/brevity?sslmode=disable"
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewPostgres(db)

	authors, err := repo.ListAllDigitalAuthors(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	generator := newArticleGenerator(client, repo)

	var wg sync.WaitGroup
	for _, author := range authors {
		wg.Go(func() {
			err := generator.generate(ctx, author)
			if err != nil {
				log.Println(err)
			}
		})
	}

	wg.Wait()
}

type articleGenerator struct {
	client      openai.Client
	repo        repository.Repository
	schemaParam openai.ResponseFormatJSONSchemaJSONSchemaParam
}

func newArticleGenerator(client openai.Client, repo repository.Repository) *articleGenerator {
	articleSchema := createJSONSchema[Article]()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "Article",
		Description: openai.String("An article / blog post that can be found on platforms like Substack"),
		Schema:      articleSchema,
		Strict:      openai.Bool(true),
	}

	return &articleGenerator{
		client:      client,
		repo:        repo,
		schemaParam: schemaParam,
	}
}

func (a *articleGenerator) generate(ctx context.Context, author *repository.DigitalAuthor) error {
	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(author.SystemPrompt),
			openai.UserMessage("Write about a random topic of your specialty"),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: a.schemaParam,
			},
		},
		Model: "minimax/minimax-m2.5",
	})
	if err != nil {
		return err
	}

	if len(chat.Choices) == 0 {
		return errors.New("error empty chat completions")
	}
	var output Article
	if err := json.Unmarshal([]byte(chat.Choices[0].Message.Content), &output); err != nil {
		return err
	}

	err = a.repo.CreateArticle(ctx, &repository.Article{
		Slug:        output.Slug,
		Content:     output.Content,
		Title:       output.Title,
		Description: output.Description,
		AuthorID:    author.ID,
	})

	return err
}

func createJSONSchema[T any]() any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: true,
		DoNotReference:            true,
	}

	var v T
	schema := reflector.Reflect(v)

	return schema
}
