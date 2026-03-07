package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sup25/gobuy/internal/auth/models"
	"github.com/sup25/gobuy/internal/auth/service"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest

		// JSON binding error
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError("Invalid input data", http.StatusBadRequest))
			return
		}

		// Validation error
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			utils.RespondWithError(c, utils.NewAppErrorWithDetails(
				"Validation failed",
				http.StatusBadRequest,
				utils.FormatValidationErrors(err),
			))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Service error (already returns AppError)
		userResp, err := service.CreateUserService(ctx, req, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusCreated, "User created successfully", userResp)
	}
}

func VerifyEmailController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		fmt.Println("Token:", token)

		if token == "" {
			utils.RespondWithError(c, utils.NewAppError("verification token required", 400))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := service.VerifyEmailService(ctx, token, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusOK, "Email verified successfully", nil)
	}
}

func ResendVerificationEmailController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError("invalid request", 400))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := service.ResendVerificationEmailService(ctx, req.Email, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, 200, "Verification email sent", nil)
	}
}

func LoginUserController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest

		// JSON binding error
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError("Invalid JSON format", http.StatusBadRequest))
			return
		}

		// Validation error
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			utils.RespondWithError(c, utils.NewAppErrorWithDetails(
				"Validation failed",
				http.StatusBadRequest,
				utils.FormatValidationErrors(err),
			))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Service error (already returns AppError)
		accessToken, refreshToken, err := service.LoginUserService(ctx, req, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}
		secure := os.Getenv("ENV") == "production"
		c.SetCookie("access_token", accessToken, 15*60, "/", "", secure, true)
		c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", secure, true)

		utils.SendSuccess(c, http.StatusOK, "User logged in successfully", gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func GoogleLoginController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Code string `json:"code" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError("code is required", 400))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Exchange code for ID token
		idToken, err := utils.ExchangeCodeForIDToken(req.Code)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		loginResp, err := service.GoogleLoginService(ctx, idToken, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, 200, "Login successful", loginResp)
	}
}

func ForgotPasswordController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError("invalid request", 400))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := service.ForgotPasswordService(ctx, req.Email, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, 200, "Password reset email sent", nil)
	}
}

func ResetPasswordController(client *mongo.Client) gin.HandlerFunc {
	{
		return func(c *gin.Context) {
			var req struct {
				Password string `json:"password"`
			}
			token := c.Query("token")
			if err := c.BindJSON(&req); err != nil {
				utils.RespondWithError(c, utils.NewAppError("invalid request", 400))
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()

			err := service.ResetPasswordService(ctx, token, req.Password, client)
			if err != nil {
				utils.RespondWithError(c, err)
				return
			}

			utils.SendSuccess(c, 200, "Password reset successfully", nil)
		}
	}
}

func RefreshTokenController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			utils.RespondWithError(c,
				utils.NewAppError("No refresh token", 401))
			return
		}
		fmt.Println("Received refresh token:", refreshToken)
		access, newRefresh, expiry, err :=
			service.RefreshTokenService(
				context.Background(),
				client,
				refreshToken,
			)

		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		secure := os.Getenv("ENV") == "production"

		c.SetCookie("access_token", access, 900, "/", "", secure, true)

		c.SetCookie(
			"refresh_token",
			newRefresh,
			int(time.Until(expiry).Seconds()),
			"/",
			"",
			secure,
			true,
		)

		utils.SendSuccess(c, 200, "Token refreshed", nil)
	}
}

func LogoutController() gin.HandlerFunc {
	return func(c *gin.Context) {
		secure := os.Getenv("ENV") == "production"
		c.SetCookie("access_token", "", -1, "/", "", secure, true)
		c.SetCookie("refresh_token", "", -1, "/", "", secure, true)
		utils.SendSuccess(c, 200, "Logout successful", nil)
	}
}
