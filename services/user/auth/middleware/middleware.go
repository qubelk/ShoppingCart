package middleware

import (
	"net/http"
	"user/auth"

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

		login, err := auth.ValidateJWT(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid authorization token",
			})
			return
		}

		ctx.Set("login", login)
		ctx.Next()
	}
}
