package server

import (
	"github.com/jmoiron/sqlx"
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/service"
)

func InitializeArticleController(db *sqlx.DB) *controller.ArticleController {
	articleRepository := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepository)
	articleController := controller.NewArticleController(articleService)
	return articleController
}

func InitializeAuthService(db *sqlx.DB, tokenSecret string) service.AuthService {
	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, tokenSecret)
	return authService
}

func InitializeLLMAPIKeyController(db *sqlx.DB, crypter service.Crypter) *controller.LLMAPIKeyController {
	llmapiKeyRepository := repository.NewLLMAPIKeyRepository(db)
	llmapiKeyService := service.NewLLMAPIKeyService(llmapiKeyRepository, crypter)
	llmapiKeyController := controller.NewLLMAPIKeyController(llmapiKeyService)
	return llmapiKeyController
}
