package service

import (
	"context"
	"time"

	"github.com/sup25/gobuy/config/db"
	userModels "github.com/sup25/gobuy/internal/user/models"

	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

func ChangePasswordService(ctx context.Context, userID, oldPassword, newPassword string, client *mongo.Client) error {
	userCollection := db.OpenCollection("users", client)

	// Convert string userID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return utils.NewAppError("invalid user ID", 400)
	}

	// Get the existing user and verify old password in one DB call
	var user userModels.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.NewAppError("user not found", 404)
		}
		return utils.NewAppError("database error", 500)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return utils.NewAppError("old password is incorrect", 400)
	}

	if oldPassword == newPassword {
		return utils.NewAppError("new password cannot be the same as old password", 400)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewAppError("failed to hash password", 500)
	}

	// Update password in DB
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"password":   string(hashedPassword),
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return utils.NewAppError("database error", 500)
	}

	return nil
}

/* GetUserFull returns the full user object (use for internal donot expose in API) */
func GetUserFull(ctx context.Context, userID string, client *mongo.Client) (*userModels.User, error) {
	userCollection := db.OpenCollection("users", client)

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, utils.NewAppError("invalid user ID", 400)
	}

	var user userModels.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.NewAppError("user not found", 404)
		}
		return nil, utils.NewAppError("database error", 500)
	}

	return &user, nil
}
