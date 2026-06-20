package main

import (
	"cart/internal/user/handler"
	"cart/internal/user/repository"
	"cart/internal/user/service"
	"cart/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewDataBase(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db)
	serv := service.New(repo)
	hand := handler.New(serv)

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.File("index.html")
	})

	r.POST("/register", hand.Register)
	r.POST("/login", hand.Login)

	profile := r.Group("/profile")
	profile.Use(middleware.AuthMiddleware(serv))
	{
		profile.GET("/", hand.GetProfile)
		profile.DELETE("/", hand.Delete)
	}

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
