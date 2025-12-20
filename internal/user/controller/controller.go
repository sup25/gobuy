package userController

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	userService "github.com/sup25/gobuy/internal/user/service"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserProfileController returns a gin.HandlerFunc
func GetUserProfileController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract userID from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			utils.RespondWithError(c, utils.NewAppError("unauthenticated", http.StatusUnauthorized))
			return
		}

		// Use context with timeout for DB operation
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Call the service function
		profile, err := userService.GetUserProfileService(ctx, userID.(string), client)
		if err != nil {
			utils.RespondWithError(c, err) // service already returns AppError
			return
		}

		// Send success response
		utils.SendSuccess(c, http.StatusOK, "User profile fetched successfully", gin.H{
			"user": profile,
		})
	}
}
