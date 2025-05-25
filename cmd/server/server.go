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
	"github.com/tuananhlai/brevity-go/internal/otelsdk"
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

	articleController := InitializeArticleController(db)
	authController := InitializeAuthController(db)

	// == Otel Setup ==
	otelShutdown, err := otelsdk.Setup(globalCtx, otelsdk.SetupConfig{
		Mode:             cfg.Mode,
		ServiceName:      serviceName,
		CollectorGrpcURL: cfg.Otel.CollectorGrpcURL,
	})
	if err != nil {
		log.Fatalf("error initializing opentelemetry sdk: %s", err)
	}

	// == Gin Setup ==
	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	// TODO: Reconfigure CORs before production deployment.
	r.Use(cors.Default())

	// == Routes ==
	articleController.RegisterRoutes(r)
	authController.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Server.Port),
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
