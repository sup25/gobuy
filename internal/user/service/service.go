package service

import (
	"context"

	"github.com/sup25/gobuy/config/db"
	userModels "github.com/sup25/gobuy/internal/user/models"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserProfileService(ctx context.Context, userID string, client *mongo.Client) (userModels.UserResponse, error) {
	userCollection := db.OpenCollection("users", client)

	// Convert string ID to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return userModels.UserResponse{}, utils.NewAppError("invalid user ID", 400)
	}

	var user userModels.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return userModels.UserResponse{}, utils.NewAppError("user not found", 404)
		}
		return userModels.UserResponse{}, utils.NewAppError("database error", 500)
	}

	// Convert to response DTO
	return user.ToResponse(), nil
}
