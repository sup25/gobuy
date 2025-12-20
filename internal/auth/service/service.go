package service

import (
	"context"
	"time"

	"github.com/sup25/gobuy/config/db"
	"github.com/sup25/gobuy/internal/auth/models"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserService(ctx context.Context, req models.RegisterRequest, client *mongo.Client) (models.UserResponse, error) {
	userCollection := db.OpenCollection("users", client)

	// Check duplicate email
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return models.UserResponse{}, utils.NewAppError("email already registered", 400)
	}
	if err != mongo.ErrNoDocuments {
		return models.UserResponse{}, utils.NewAppError("database error", 500)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, utils.NewAppError("failed to hash password", 500)
	}

	// Prepare user
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "USER",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user
	res, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return models.UserResponse{}, utils.NewAppError("failed to create user", 500)
	}

	// Assign MongoDB _id
	user.ID = res.InsertedID.(primitive.ObjectID)

	// Return safe user info
	return user.ToResponse(), nil
}

func LoginUserService(
	ctx context.Context,
	req models.LoginRequest,
	client *mongo.Client,
) (string, string, error) {

	userCollection := db.OpenCollection("users", client)

	// Find user
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{
		"email": req.Email,
	}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return "", "", utils.NewAppError("Invalid email or password", 400)
	}
	if err != nil {
		return "", "", utils.NewAppError("database error", 500)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", "", utils.NewAppError("Invalid email or password", 400)
	}

	// Generate tokens
	accessToken, refreshToken, refreshExpiry, err := utils.GenerateAuthTokens(user)
	if err != nil {
		return "", "", utils.NewAppError("Failed to generate tokens", 500)
	}

	// Hash the refresh token before storing
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", utils.NewAppError("Failed to hash token", 500)
	}

	// Store refresh token in database
	refreshTokenDoc := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: string(hashedToken),
		ExpiresAt: refreshExpiry,
		CreatedAt: time.Now(),
		IsRevoked: false,
	}

	refreshTokenCollection := db.OpenCollection("refresh_tokens", client)
	_, err = refreshTokenCollection.InsertOne(ctx, refreshTokenDoc)
	if err != nil {
		return "", "", utils.NewAppError("Failed to store refresh token", 500)
	}

	return accessToken, refreshToken, nil
}
