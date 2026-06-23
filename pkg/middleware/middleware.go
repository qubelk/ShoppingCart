package middleware

import (
	"net/http"
	"pkg/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := auth.New()
		token, err := ctx.Cookie("jwt-token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "need authorization first",
			})
			return
		}

		id, login, err := auth.ValidateJWT(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid authorization token",
			})
			return
		}

		ctx.Set("user_id", id)
		ctx.Set("login", login)
		ctx.Next()
	}
}
