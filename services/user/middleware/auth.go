package middleware

import (
	"net/http"
	"user/internal/user/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("jwt-token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "need authorization first",
			})
			return
		}

		id, err := userService.ValidateJWT(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid authorization token",
			})
			return
		}

		ctx.Set("id", id)
		ctx.Next()
	}
}
