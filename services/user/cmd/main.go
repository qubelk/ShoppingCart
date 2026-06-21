package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/user"
	"user/internal/user/handler"
	"user/internal/user/repository"
	"user/internal/user/service"
	"user/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	user.LogInfo("Log")

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
		ctx.File("./index.html")
	})

	users := r.Group("/users")
	{
		users.POST("", hand.Register)
		users.POST("/auth", hand.Login)
	}

	users.Use(middleware.AuthMiddleware(serv))
	{
		users.GET("/:id", hand.GetProfile)
		users.DELETE("/:id", hand.Delete)
	}

	srv := http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
