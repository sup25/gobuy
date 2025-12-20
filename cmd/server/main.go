package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sup25/gobuy/config/db"

	privateRoutes "github.com/sup25/gobuy/routes/private"
	publicRoutes "github.com/sup25/gobuy/routes/public"
)

func main() {
	router := gin.Default()

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
