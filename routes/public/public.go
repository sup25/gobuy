package public

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sup25/gobuy/internal/auth/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

func PublicRoute(router *gin.Engine, client *mongo.Client) {
	router.POST("/register", controller.RegisterUserController(client))
	router.POST("/login", controller.LoginUserController(client))
	router.GET("/verify-email", controller.VerifyEmailController(client))
	router.POST("/resend-verification-email", controller.ResendVerificationEmailController(client))
	router.POST("/forgot-password", controller.ForgotPasswordController(client))
	router.POST("/reset-password", controller.ResetPasswordController(client))
	router.POST("/google-login", controller.GoogleLoginController(client))
}
