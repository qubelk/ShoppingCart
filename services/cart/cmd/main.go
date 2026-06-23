package main

import (
	"cart/internal/handler"
	"cart/internal/repository"
	"cart/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pkg/logs"
	"pkg/middleware"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logs.Init()
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	vlk, err := repository.NewValkeyClient()
	if err != nil {
		log.Fatal(err)
	}
	defer vlk.Close()

	repo := repository.New(vlk)
	serv := service.New(repo)
	hand := handler.New(serv)

	r := gin.Default()

	api := r.Group("/cart")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/", hand.GetCart)
		api.POST("/items", hand.AddItem)
		api.PUT("/items/quantity", hand.UpdateQuantity)
		api.PUT("/items/remove", hand.RemoveItem)
		api.POST("/clear", hand.CleanCart)
		api.GET("/ttl", hand.GetCartTTL)
	}

	srv := http.Server{
		Addr:    ":3003",
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
