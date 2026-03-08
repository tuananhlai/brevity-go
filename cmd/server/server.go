package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/tuananhlai/brevity-go/internal/config"
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/encryption"
	"github.com/tuananhlai/brevity-go/internal/store"
	"github.com/tuananhlai/brevity-go/internal/telemetry"
	"github.com/tuananhlai/brevity-go/internal/token"
)

const (
	// TODO: sign with private key instead
	accessTokenSecret = "secret"
)

func Run() {
	cfg := config.MustLoadConfig()

	globalCtx := context.Background()
	db := otelsqlx.MustConnect("postgres", cfg.Database.URL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL))

	s := store.NewPostgresStore(db)
	tokenIssuer := token.NewIssuer(accessTokenSecret)
	articleController := initializeArticleController(s)
	authController := initializeAuthController(s, tokenIssuer)
	healthController := controller.NewHealthController()
	encryptionService, err := encryption.New([]byte(cfg.Encryption.Key))
	if err != nil {
		log.Fatalf("error initializing encryption service: %s", err)
	}
	llmAPIKeyController := initializeLLMAPIKeyController(s, encryptionService)
	authMiddleware := controller.AuthMiddleware(tokenIssuer)

	// == Otel Setup ==
	err = telemetry.Setup(globalCtx)
	if err != nil {
		log.Fatalf("error initializing opentelemetry sdk: %s", err)
	}

	logger := telemetry.Logger("github.com/tuananhlai/brevity-go/cmd/server")

	// == Gin Setup ==
	r := gin.Default()
	r.Use(otelgin.Middleware("main-server"))
	r.Use(cors.New(getCorsConfig()))

	// == Routes ==
	r.GET("/health/liveness", healthController.CheckLiveness)
	r.POST("/v1/auth/sign-up", authController.Register)
	r.POST("/v1/auth/sign-in", authController.Login)
	r.GET("/v1/article-previews", articleController.ListPreviews)
	r.GET("/v1/articles/:slug", articleController.GetBySlug)
	r.GET("/v1/auth/me", authMiddleware, authController.GetCurrentUser)
	r.POST("/v1/llm-api-keys", authMiddleware, llmAPIKeyController.CreateLLMAPIKey)
	r.GET("/v1/llm-api-keys", authMiddleware, llmAPIKeyController.ListLLMAPIKeys)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Server.Port),
		Handler: r.Handler(),
	}

	logger.Info("Server started on port", "port", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

	// TODO: enable ReleaseMode for Gin
	// TODO: add graceful shutdown
}

var allowedOrigins = []string{
	// Local development
	"http://localhost:5173",
	// Production environments
	"https://brevity-next.vercel.app",
	"https://brevity.laituananh.com",
}

func getCorsConfig() cors.Config {
	cfg := cors.DefaultConfig()
	cfg.AllowOriginFunc = func(origin string) bool {
		if slices.Contains(allowedOrigins, origin) {
			return true
		}

		// Allow Vercel preview deployments (e.g., https://brevity-next-*.vercel.app)
		if strings.HasPrefix(origin, "https://brevity-next-") &&
			strings.HasSuffix(origin, ".vercel.app") {
			return true
		}

		// Deny all other origins
		return false
	}
	cfg.AllowCredentials = true
	return cfg
}
