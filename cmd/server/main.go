package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	if err := router.Run(":8080"); err != nil {

		panic(err)
	}
}
