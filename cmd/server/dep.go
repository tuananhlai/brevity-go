package server

import (
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/llmapikey"
	"github.com/tuananhlai/brevity-go/internal/store"
	"github.com/tuananhlai/brevity-go/internal/token"
)

func initializeArticleController(s *store.PostgresStore) *controller.ArticleController {
	return controller.NewArticleController(s)
}

func initializeAuthController(s *store.PostgresStore, tokenIssuer *token.AccessTokenIssuer) *controller.AuthController {
	return controller.NewAuthController(s, tokenIssuer)
}

func initializeLLMAPIKeyController(s *store.PostgresStore, crypter llmapikey.Crypter) *controller.LLMAPIKeyController {
	manager := llmapikey.NewManager(s, crypter)
	return controller.NewLLMAPIKeyController(manager)
}
