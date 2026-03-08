package server

import (
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/encryption"
	"github.com/tuananhlai/brevity-go/internal/llmapikey"
	"github.com/tuananhlai/brevity-go/internal/store"
	"github.com/tuananhlai/brevity-go/internal/token"
)

func initializeArticleController(s *store.Store) *controller.ArticleController {
	return controller.NewArticleController(s)
}

func initializeAuthController(s *store.Store, tokenIssuer *token.AccessTokenIssuer) *controller.AuthController {
	return controller.NewAuthController(s, tokenIssuer)
}

func initializeLLMAPIKeyController(s *store.Store, crypter *encryption.Cipher) *controller.LLMAPIKeyController {
	manager := llmapikey.NewManager(s, crypter)
	return controller.NewLLMAPIKeyController(manager)
}
