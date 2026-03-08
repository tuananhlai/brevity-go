package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuananhlai/brevity-go/internal/token"
)

// AuthMiddleware stops HTTP requests that do not have a valid access token from reaching the HTTP handler.
func AuthMiddleware(tokenIssuer *token.AccessTokenIssuer) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		accessToken, ok := extractAccessTokenFromRequest(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    CodeUnauthorized,
				Message: "access token not found",
			})
			ginCtx.Abort()
			return
		}

		userID, err := tokenIssuer.Verify(accessToken)
		if err != nil {
			ginCtx.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    CodeUnauthorized,
				Message: fmt.Sprintf("error verifying access token: %s", err.Error()),
			})
			ginCtx.Abort()
			return
		}

		setContextUserID(ginCtx, userID)
		ginCtx.Next()
	}
}
