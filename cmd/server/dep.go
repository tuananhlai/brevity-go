package server

import (
	"github.com/tuananhlai/brevity-go/internal/articles"
	"github.com/tuananhlai/brevity-go/internal/auth"
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/llmapikey"
	"github.com/tuananhlai/brevity-go/internal/repository"
)

func initializeArticleController(repo *repository.Postgres) *controller.ArticleController {
	articleService := articles.NewService(repo)
	articleController := controller.NewArticleController(articleService)
	return articleController
}

func initializeAuthService(repo *repository.Postgres, tokenSecret string) auth.Service {
	authService := auth.NewService(repo, tokenSecret)
	return authService
}

func initializeLLMAPIKeyController(repo *repository.Postgres, crypter llmapikey.Crypter) *controller.LLMAPIKeyController {
	llmapiKeyService := llmapikey.NewService(repo, crypter)
	llmapiKeyController := controller.NewLLMAPIKeyController(llmapiKeyService)
	return llmapiKeyController
}
