package users

import (
	"net/http"
	"outshort/app/common"

	"github.com/gin-gonic/gin"
)

func AuthRequired(storage *Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := common.GetAuthTokenFromHeader(context)
		user, err := storage.GetUserInfo(token)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		context.Set("user", *user)
		context.Next()
	}
}
