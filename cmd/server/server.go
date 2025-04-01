package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/tuananhlai/brevity-go/internal/config"
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/otelsdk"
	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/service"
)

const (
	shutdownTimeout = 5 * time.Second
	serviceName     = "brevity"
)

func Run() {
	cfg := config.MustLoadConfig()
	if cfg.Mode == config.ModeRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	globalCtx := context.Background()
	db := otelsqlx.MustConnect("postgres", cfg.Database.URL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL))

	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewAuthController(authService)

	// == Otel Setup ==
	otelShutdown, err := otelsdk.Setup(globalCtx, otelsdk.SetupConfig{
		Mode:             cfg.Mode,
		ServiceName:      serviceName,
		CollectorGrpcURL: cfg.Otel.CollectorGrpcURL,
	})
	if err != nil {
		log.Fatalf("error initializing opentelemetry sdk: %s", err)
	}

	logger := otelsdk.Logger("article")
	tracer := otelsdk.Tracer("article")

	// == Gin Setup ==
	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))

	r.POST("/auth/register", authController.Register)

	r.GET("/article-previews", func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "article.listPreviews")
		defer span.End()

		logger.InfoContext(ctx, "Received request to get article previews")

		articles, err := articleService.ListPreviews(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, articles)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
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

	log.Println("Shutting down server...")

	timeoutCtx, cancel := context.WithTimeout(globalCtx, shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Printf("Failed to shutdown server gracefully: %v", err)
	}

	if err := otelShutdown(timeoutCtx); err != nil {
		log.Printf("Failed to shutdown opentelemetry sdk: %v", err)
	}

	// Close the database connection after the server has been shutdown to ensure in-flight requests are completed.
	if err := db.Close(); err != nil {
		log.Printf("Failed to close database connection gracefully: %v", err)
	}

	<-timeoutCtx.Done()
	log.Println("Server shutdown complete.")
}
