package private

import (
	"github.com/gin-gonic/gin"
	userController "github.com/sup25/gobuy/internal/user/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupUserRoutes(rg *gin.RouterGroup, client *mongo.Client) {
	userRoutes := rg.Group("/user")
	{
		userRoutes.GET("/me", userController.GetUserProfileController(client))
		userRoutes.POST("/change-password", userController.ChangePasswordController(client))
	}
}
