//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"

	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/service"
)

func InitializeArticleController(db *sqlx.DB) *controller.ArticleController {
	wire.Build(controller.NewArticleController, service.NewArticleService, repository.NewArticleRepository)

	return &controller.ArticleController{}
}

func InitializeAuthController(db *sqlx.DB, opts ...service.AuthServiceOption) *controller.AuthController {
	wire.Build(controller.NewAuthController, service.NewAuthService, repository.NewAuthRepository)

	return &controller.AuthController{}
}
