package private

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sup25/gobuy/internal/user/controller"
	"github.com/sup25/gobuy/pkg/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

// PrivateRoute registers private routes
func PrivateRoute(router *gin.Engine, client *mongo.Client) {
	// Create a private route group
	private := router.Group("/private")
	private.Use(middleware.AuthMiddleware()) // attach auth middleware

	// Register routes
	private.GET("/profile", controller.GetUserProfileController(client))
}
