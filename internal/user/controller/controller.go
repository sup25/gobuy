package userController

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sup25/gobuy/internal/user/models"
	userService "github.com/sup25/gobuy/internal/user/service"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserProfileController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract userID from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			utils.RespondWithError(c, utils.NewAppError("unauthenticated", http.StatusUnauthorized))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		profile, err := userService.GetUserProfileService(ctx, userID.(string), client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusOK, "User profile fetched successfully", gin.H{
			"user": profile,
		})
	}
}

func ChangePasswordController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.ChangePasswordRequest
		fmt.Printf("Request body: %+v\n", req)

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError("invalid request", http.StatusBadRequest))
			return
		}

		// Extract userID from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			utils.RespondWithError(c, utils.NewAppError("unauthenticated", http.StatusUnauthorized))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := userService.ChangePasswordService(ctx, userID.(string), req.OldPassword, req.NewPassword, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusOK, "Password changed successfully", nil)
	}
}
