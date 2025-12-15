package server

import (
	"github.com/jmoiron/sqlx"

	"github.com/tuananhlai/brevity-go/internal/articles"
	"github.com/tuananhlai/brevity-go/internal/auth"
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/llmapikey"
)

func initializeArticleController(db *sqlx.DB) *controller.ArticleController {
	articleRepository := articles.NewRepository(db)
	articleService := articles.NewService(articleRepository)
	articleController := controller.NewArticleController(articleService)
	return articleController
}

func initializeAuthService(db *sqlx.DB, tokenSecret string) auth.Service {
	authRepository := auth.NewRepository(db)
	authService := auth.NewService(authRepository, tokenSecret)
	return authService
}

func initializeLLMAPIKeyController(db *sqlx.DB, crypter llmapikey.Crypter) *controller.LLMAPIKeyController {
	llmapiKeyRepository := llmapikey.NewRepository(db)
	llmapiKeyService := llmapikey.NewService(llmapiKeyRepository, crypter)
	llmapiKeyController := controller.NewLLMAPIKeyController(llmapiKeyService)
	return llmapiKeyController
}
