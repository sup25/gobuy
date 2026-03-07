package private

import (
	"github.com/gin-gonic/gin"
	"github.com/sup25/gobuy/pkg/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func PrivateRoute(router *gin.Engine, client *mongo.Client) {
	private := router.Group("/api/v1")
	private.Use(middleware.AuthMiddleware())

	// Setup all private routes
	SetupUserRoutes(private, client)
	SetupProductRoutes(private, client)
}
