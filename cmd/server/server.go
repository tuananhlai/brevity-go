package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/tuananhlai/brevity-go/internal/config"
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/controller/middleware"
	"github.com/tuananhlai/brevity-go/internal/encryption"
	"github.com/tuananhlai/brevity-go/internal/otelsdk"
)

const (
	shutdownTimeout = 5 * time.Second
)

func Run() {
	cfg := config.MustLoadConfig()
	if cfg.Mode == config.ModeRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	globalCtx := context.Background()
	db := otelsqlx.MustConnect("postgres", cfg.Database.URL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL))

	articleController := initializeArticleController(db)
	authService := initializeAuthService(db, cfg.Auth.TokenSecret)
	authController := controller.NewAuthController(authService)
	healthController := controller.NewHealthController()
	encryptionService, err := encryption.New([]byte(cfg.Encryption.Key))
	if err != nil {
		log.Fatalf("error initializing encryption service: %s", err)
	}
	llmAPIKeyController := initializeLLMAPIKeyController(db, encryptionService)
	authMiddleware := middleware.AuthMiddleware(authService)

	// == Otel Setup ==
	otelShutdown, err := otelsdk.Setup(globalCtx, otelsdk.SetupConfig{
		Debug: cfg.Mode == config.ModeDev,
	})
	if err != nil {
		log.Fatalf("error initializing opentelemetry sdk: %s", err)
	}

	logger := otelsdk.Logger("github.com/tuananhlai/brevity-go/cmd/server")

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

	go func() {
		logger.Info("Server started on port", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	// == Graceful Shutdown ==
	if cfg.Mode == config.ModeDev {
		log.Println("Server is running in dev mode, skipping graceful shutdown.")
		return
	}

	logger.Info("Shutting down server...")

	timeoutCtx, cancel := context.WithTimeout(globalCtx, shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		logger.Warn("Failed to shutdown server gracefully", "error", err)
	}

	if err := otelShutdown(timeoutCtx); err != nil {
		logger.Warn("Failed to shutdown opentelemetry sdk", "error", err)
	}

	// Close the database connection after the server has been shutdown to ensure in-flight requests are completed.
	if err := db.Close(); err != nil {
		logger.Warn("Failed to close database connection gracefully", "error", err)
	}

	<-timeoutCtx.Done()
	logger.Info("Server shutdown complete.")
}

func getCorsConfig() cors.Config {
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{"http://localhost:3000", "https://brevity-next.vercel.app/"}
	cfg.AllowCredentials = true
	return cfg
}
