package main

import (
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://6dvt25.csb.app"}, // frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"project": "gobuy-server",
			"status":  "running successfully!",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	mongoClient := db.Client

	// Make sure function is exported (capital P)
	publicRoutes.PublicRoute(router, mongoClient)
	privateRoutes.PrivateRoute(router, mongoClient)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
