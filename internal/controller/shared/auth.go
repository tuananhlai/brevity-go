package shared

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookieName   = "$access_token"
	accessTokenCookiePath   = "/"
	accessTokenCookieMaxAge = 60 * 60 * 24 * 30 // 30 days
	contextKeyUserID        = "$userID"
)

func SetContextUserID(ginCtx *gin.Context, userID string) {
	ginCtx.Set(contextKeyUserID, userID)
}

// GetContextUserID extracts the current user ID from the Gin context.
func GetContextUserID(ginCtx *gin.Context) (string, error) {
	rawUserID, ok := ginCtx.Get(contextKeyUserID)
	if !ok {
		return "", errors.New("userID not found in context")
	}

	userID, ok := rawUserID.(string)
	if !ok {
		return "", errors.New("userID is not a string")
	}

	return userID, nil
}

// ExtractAccessTokenFromRequest returns the access token from the HTTP request. If the access token is found,
// the second return value will be true, otherwise it will be false.
func ExtractAccessTokenFromRequest(ginCtx *gin.Context) (string, bool) {
	token, err := ginCtx.Cookie(accessTokenCookieName)
	if err != nil {
		return "", false
	}

	return token, true
}

func SetAccessTokenCookie(ginCtx *gin.Context, token string) {
	ginCtx.SetCookie(accessTokenCookieName, token, accessTokenCookieMaxAge, accessTokenCookiePath, "", true, true)
}
