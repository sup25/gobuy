package controller

import (
	"context"
	"net/http"
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

		c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", false, true)

		utils.SendSuccess(c, http.StatusOK, "User logged in successfully", gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}
