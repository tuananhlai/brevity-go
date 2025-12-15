package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tuananhlai/brevity-go/internal/auth"
	"github.com/tuananhlai/brevity-go/internal/controller/shared"
)

// AuthMiddleware stops HTTP requests that do not have a valid access token from reaching the HTTP handler.
func AuthMiddleware(authService auth.Service) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		token, ok := shared.ExtractAccessTokenFromRequest(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, shared.ErrorResponse{
				Code:    shared.CodeUnauthorized,
				Message: "access token not found",
			})
			ginCtx.Abort()
			return
		}

		userID, err := authService.VerifyAccessToken(ginCtx.Request.Context(), token)
		if err != nil {
			ginCtx.JSON(http.StatusUnauthorized, shared.ErrorResponse{
				Code:    shared.CodeUnauthorized,
				Message: fmt.Sprintf("error verifying access token: %s", err.Error()),
			})
			ginCtx.Abort()
			return
		}

		shared.SetContextUserID(ginCtx, userID)
		ginCtx.Next()
	}
}
