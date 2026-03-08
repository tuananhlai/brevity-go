package jobs

import (
	"context"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/tuananhlai/brevity-go/internal/config"
	"github.com/tuananhlai/brevity-go/internal/genarticle"
	"github.com/tuananhlai/brevity-go/internal/store"
)

func RunGenerateArticle() {
	cfg := config.MustLoadConfig()

	ctx := context.Background()
	baseURL := "https://openrouter.ai/api/v1"
	client := openai.NewClient(option.WithBaseURL(baseURL), option.WithAPIKey(cfg.LLMAPIKey))

	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalln(err)
	}

	s := store.New(db)

	authors, err := s.ListDigitalAuthorsWithArticleSlugs(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	generator := genarticle.New(client)

	var wg sync.WaitGroup
	for _, author := range authors {
		wg.Go(func() {
			result, err := generator.Generate(ctx, author.SystemPrompt, author.ArticleSlugs)
			if err != nil {
				log.Printf("generation for author %s failed: %v\n", author.ID, err)
				return
			}

			err = s.CreateArticle(ctx, &store.Article{
				Slug:        result.Slug,
				Title:       result.Title,
				Description: result.Description,
				Content:     result.Content,
				AuthorID:    author.ID,
			})
			log.Printf("generation for author %s completed. error = %v\n", author.ID, err)
		})
	}

	wg.Wait()
}
