package main

import (
	"log"
	"os"

	"github.com/MosinFAM/subs-app/internal/db"
	"github.com/MosinFAM/subs-app/internal/handlers"
	"github.com/MosinFAM/subs-app/internal/middleware"
	"github.com/MosinFAM/subs-app/internal/repo"
	"github.com/gin-gonic/gin"

	_ "github.com/MosinFAM/subs-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Marketplace API
// @version 1.0
// @description REST API for a marketplace with user auth and ads

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	repo := repo.NewPostgresRepo(conn)
	h := &handlers.Handler{Repo: repo}

	r := gin.Default()
	r.Use(middleware.GinLogger())

	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.POST("", h.CreateSubscription)
		subscriptions.GET("", h.ListSubscriptions)
		subscriptions.GET(":id", h.GetSubscription)
		subscriptions.PUT(":id", h.UpdateSubscription)
		subscriptions.DELETE(":id", h.DeleteSubscription)
	}

	r.GET("/summary", h.GetSubscriptionSummary)

	// Swagger docs only in non-prod
	if os.Getenv("ENV") != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		os.Exit(1)
	}
}
