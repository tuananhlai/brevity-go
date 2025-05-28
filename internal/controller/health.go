package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) RegisterRoutes(router *gin.Engine) {
	router.GET("/health/liveness", c.CheckLiveness)
}

func (c *HealthController) CheckLiveness(ginCtx *gin.Context) {
	ginCtx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
