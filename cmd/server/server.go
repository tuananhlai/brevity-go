package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/tuananhlai/brevity-go/internal/config"
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

	db := sqlx.MustConnect("postgres", cfg.Database.URL)
	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)

	// == Otel Setup ==
	resource, err := newResource(serviceName)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	log.Println("Resource", resource)

	// == Gin Setup ==
	r := gin.Default()

	r.GET("/article-previews", func(c *gin.Context) {
		articles, err := articleService.ListPreviews(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown server gracefully: %v", err)
	}

	// Close the database connection after the server has been shutdown to ensure in-flight requests are completed.
	if err := db.Close(); err != nil {
		log.Printf("Failed to close database connection gracefully: %v", err)
	}

	<-ctx.Done()
	log.Println("Server shutdown complete.")
}

func newResource(serviceName string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	))
}
