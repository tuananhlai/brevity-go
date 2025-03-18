package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tuananhlai/brevity-go/internal/article"
	"github.com/tuananhlai/brevity-go/internal/config"
)

func Run() {
	cfg := config.MustLoadConfig()
	db := sqlx.MustConnect("postgres", cfg.Database.URL)
	articleRepo := article.NewRepository(db)
	articleService := article.NewService(articleRepo)

	r := gin.Default()

	r.GET("/article-previews", func(c *gin.Context) {
		articles, err := articleService.ListPreviews(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, articles)
	})

	r.Run()
}
