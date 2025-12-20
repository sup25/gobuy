package public

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sup25/gobuy/internal/auth/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

func PublicRoute(router *gin.Engine, client *mongo.Client) {
	router.POST("/register", controller.RegisterUserController(client))
	router.POST("/login", controller.LoginUserController(client))
}
