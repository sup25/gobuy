package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sup25/gobuy/config/db"

	privateRoutes "github.com/sup25/gobuy/routes/private"
	publicRoutes "github.com/sup25/gobuy/routes/public"
)

func main() {
	router := gin.Default()

	// CORS setup
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "https://6dvt25.csb.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"project": "gobuy-server",
			"status":  "running successfully!",
		})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Mongo client
	mongoClient := db.Client

	// Seed admin and manager
	if err := db.SeedAdminAndManager(mongoClient); err != nil {
		log.Fatal("Failed to seed admin/manager:", err)
	}

	// Create indexes
	if err := db.CreateIndexes(mongoClient); err != nil {
		log.Fatal("Index creation failed:", err)
	}

	// Register routes
	publicRoutes.PublicRoute(router, mongoClient)
	privateRoutes.PrivateRoute(router, mongoClient)

	// Start server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
