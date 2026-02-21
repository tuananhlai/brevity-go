package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/invopop/jsonschema"
	"github.com/jmoiron/sqlx"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/tuananhlai/brevity-go/internal/repository"
)

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

	authors, err := repo.ListDigitalAuthorsWithArticleSlugs(ctx)
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

func (a *articleGenerator) generate(ctx context.Context, author *repository.DigitalAuthorWithArticleSlugs) error {
	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(author.SystemPrompt),
			openai.SystemMessage(writingStyleSystemPrompt),
			openai.SystemMessage(fmt.Sprintf("You have already written articles with the following slugs: %v. Do not write about the same topic.", author.ArticleSlugs)),
			openai.UserMessage("Write about a random topic of your specialty"),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: a.schemaParam,
			},
		},
		Model: "moonshotai/kimi-k2.5",
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

const (
	writingStyleSystemPrompt = `
Follow these requirements for clear, effective writing:

A. Audience & Document Purpose
- Identify and implicitly target a specific audience.
- Tailor explanations to what the audience knows and needs to learn.
- Begin by stating document scope, audience, and key points.

B. Word & Sentence Level
- Define new or unfamiliar terms before use.
- Use terms consistently across the document.
- Avoid ambiguous pronouns with unclear referents.
- Prefer active voice when it improves clarity.
- Choose specific nouns and strong verbs.
- Use punctuation correctly to aid readability.
- Focus each sentence on a single idea.
- Eliminate unnecessary words.

C. Lists, Tables, and Parallelism
- Convert long or complex sentences into lists when useful.
- Use numbered lists for ordered steps; start numbered items with imperative verbs.
- Use bullet lists when order is not important.
- Keep list items parallel in structure and grammar.
- Introduce lists and tables with a clear lead-in sentence.

D. Paragraph & Section Structure
- Begin paragraphs with strong topic sentences.
- Keep each paragraph focused on a single topic.
- Break long topics into clear sections with descriptive headings.

E. Document Organization
- Follow a consistent style guide or template.
- Think from the reader’s perspective when choosing structure.
- Outline large documents before writing or reorganize after drafting.
- Prefer task-oriented, actionable section headings.
- Use progressive disclosure: introduce simpler concepts before complex ones.
- Define scope early and avoid off-scope digressions.

F. Visuals & Illustrations
- Write captions that explain the key takeaway.
- Limit the amount of information in a single diagram.
- Use visual emphasis to direct reader attention.
- Ensure visuals support, not duplicate, the text.

G. Code Samples (If Applicable)
- Provide concise, correct, and runnable examples.
- Keep examples minimal but complete.
- Write short, meaningful comments; avoid commenting obvious code.
- Include both examples and, when helpful, counterexamples.
- Demonstrate increasing complexity when needed.

H. Revision & Quality Control
- Review for clarity and logical flow.
- Remove ambiguity and tighten language.
- Check consistency of terminology.
- Verify technical accuracy.
- Revise after distance or external feedback.
- Treat revision as iterative and continuous.

I. Output Requirements
- Structured with clear headings.
- Neutral, precise tone.
- Optimized for readability and scanability.
- No unnecessary verbosity.`
)
