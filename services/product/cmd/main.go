package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pkg/middleware"
	"product/internal/handler"
	"product/internal/product"
	"product/internal/repository"
	"product/internal/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	product.LogInfo("Log")

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

	products := r.Group("/products")
	products.Use(middleware.AuthMiddleware())
	{
		products.POST("/", hand.Create)
		products.GET("/:id", hand.GetProduct)
		products.DELETE("/:id", hand.Delete)
		products.GET("/search", hand.SearchProducts)
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
