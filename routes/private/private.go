package private

import (
	"github.com/gin-gonic/gin"
	userController "github.com/sup25/gobuy/internal/user/controller"
	"github.com/sup25/gobuy/pkg/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func PrivateRoute(router *gin.Engine, client *mongo.Client) {
	// Create a private route group
	private := router.Group("/api/v1")
	private.Use(middleware.AuthMiddleware()) // attach auth middleware

	private.GET("/profile", userController.GetUserProfileController(client))
	private.POST("/change-password", userController.ChangePasswordController(client))

}
