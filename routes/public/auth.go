package public

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sup25/gobuy/internal/auth/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthPublicRoutes(r *gin.RouterGroup, client *mongo.Client) {
	r.POST("/auth/register", controller.RegisterUserController(client))
	r.POST("/auth/login", controller.LoginUserController(client))
	r.POST("/auth/refresh", controller.RefreshTokenController(client))
	r.POST("/auth/verify-email", controller.VerifyEmailController(client))
	r.POST("/auth/resend-verification-email", controller.ResendVerificationEmailController(client))
	r.POST("/auth/forgot-password", controller.ForgotPasswordController(client))
	r.POST("/auth/reset-password", controller.ResetPasswordController(client))
	r.POST("/auth/google-login", controller.GoogleLoginController(client))
	r.POST("/auth/logout", controller.LogoutController())
}
