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
	"go.opentelemetry.io/contrib/bridges/otelslog"

	"github.com/tuananhlai/brevity-go/internal/config"
	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/service"
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
	db := sqlx.MustConnect("postgres", cfg.Database.URL)
	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)

	// == Otel Setup ==
	otelShutdown, err := setupOTelSDK(globalCtx)
	if err != nil {
		log.Fatalf("error initializing opentelemetry sdk: %s", err)
	}

	otelLogger := otelslog.NewLogger(serviceName)

	// == Gin Setup ==
	r := gin.Default()

	r.GET("/article-previews", func(c *gin.Context) {
		otelLogger.Info("Received request to get article previews")

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
