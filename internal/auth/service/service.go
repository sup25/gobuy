package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sup25/gobuy/config/db"
	authModels "github.com/sup25/gobuy/internal/auth/models"
	userModels "github.com/sup25/gobuy/internal/user/models"
	"github.com/sup25/gobuy/pkg/utils"
	"github.com/sup25/gobuy/pkg/utils/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserService(
	ctx context.Context,
	req authModels.RegisterRequest,
	client *mongo.Client,
) (userModels.UserResponse, error) {

	userCollection := db.OpenCollection("users", client)

	// 1️⃣ Check duplicate email
	var existingUser userModels.User
	err := userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return userModels.UserResponse{}, utils.NewAppError("email already registered", 400)
	}
	if err != mongo.ErrNoDocuments {
		return userModels.UserResponse{}, utils.NewAppError("database error", 500)
	}

	// 2️⃣ Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return userModels.UserResponse{}, utils.NewAppError("failed to hash password", 500)
	}

	// 3️⃣ Generate email verification token
	verificationToken, err := utils.GenerateSecureRandomToken(32)
	if err != nil {
		return userModels.UserResponse{}, utils.NewAppError("failed to generate verification token", 500)
	}

	verificationExpiry := time.Now().Add(24 * time.Hour)

	// 4️⃣ Prepare user
	user := userModels.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "USER",

		IsEmailVerified: false,

		EmailVerificationToken:   verificationToken,
		EmailVerificationExpires: verificationExpiry,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 5️⃣ Insert user
	res, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return userModels.UserResponse{}, utils.NewAppError("failed to create user", 500)
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	go mail.SendVerificationEmail(user.Email, verificationToken)

	return user.ToResponse(), nil
}

func VerifyEmailService(ctx context.Context, token string, client *mongo.Client) error {
	userCollection := db.OpenCollection("users", client)

	// Filter user by token and make sure token is not expired
	filter := bson.M{
		"email_verification_token": token,
		"email_verification_expires": bson.M{
			"$gt": time.Now(), // token still valid
		},
	}

	// Update user: mark email verified and remove token
	update := bson.M{
		"$set": bson.M{
			"is_email_verified": true,
			"updated_at":        time.Now(),
		},
		"$unset": bson.M{
			"email_verification_token":   "",
			"email_verification_expires": "",
		},
	}

	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return utils.NewAppError("database error", 500)
	}

	if result.MatchedCount == 0 {
		return utils.NewAppError("invalid or expired verification token", 400)
	}

	return nil
}

func ResendVerificationEmailService(ctx context.Context, email string, client *mongo.Client) error {
	userCollection := db.OpenCollection("users", client)

	var user userModels.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.NewAppError("user not found", 404)
		}
		return utils.NewAppError("database error", 500)
	}

	if user.IsEmailVerified {
		return utils.NewAppError("email already verified", 400)
	}

	// Generate new token and expiry
	token, err := utils.GenerateSecureRandomToken(32)
	if err != nil {
		return utils.NewAppError("failed to generate verification token", 500)
	}

	expiry := time.Now().Add(24 * time.Hour)

	// Update user
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"email_verification_token":   token,
			"email_verification_expires": expiry,
			"updated_at":                 time.Now(),
		},
	})
	if err != nil {
		return utils.NewAppError("database error", 500)
	}

	// Send email
	go mail.SendVerificationEmail(user.Email, token)

	return nil
}

func ForgotPasswordService(ctx context.Context, email string, client *mongo.Client) error {
	userCollection := db.OpenCollection("users", client)

	var user userModels.User
	if err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.NewAppError("user not found", 404)
		}
		return utils.NewAppError("database error", 500)
	}

	// Generate reset token
	token, err := utils.GenerateSecureRandomToken(32)
	if err != nil {
		return utils.NewAppError("failed to generate token", 500)
	}

	expiry := time.Now().Add(time.Hour) // 1 hour expiry

	// Update user
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"password_reset_token":   token,
			"password_reset_expires": expiry,
			"updated_at":             time.Now(),
		},
	})
	if err != nil {
		return utils.NewAppError("database error", 500)
	}

	// Send reset email
	go mail.SendPasswordResetEmail(user.Email, token)
	return nil
}

func ResetPasswordService(ctx context.Context, token, newPassword string, client *mongo.Client) error {
	userCollection := db.OpenCollection("users", client)

	filter := bson.M{
		"password_reset_token": token,
		"password_reset_expires": bson.M{
			"$gt": time.Now(),
		},
	}

	var user userModels.User
	if err := userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return utils.NewAppError("invalid or expired token", 400)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewAppError("failed to hash password", 500)
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"password":               string(hashedPassword),
			"password_reset_token":   "",
			"password_reset_expires": nil,
			"updated_at":             time.Now(),
		},
	})
	if err != nil {
		return utils.NewAppError("database error", 500)
	}

	return nil
}

func LoginUserService(
	ctx context.Context,
	req authModels.LoginRequest,
	client *mongo.Client,
) (string, string, error) {

	userCollection := db.OpenCollection("users", client)

	// Find user
	var user userModels.User
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

	if !user.IsEmailVerified {
		return "", "", utils.NewAppError("email not verified. please check your inbox", 401)
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
	refreshTokenDoc := authModels.RefreshToken{
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
func GoogleLoginService(ctx context.Context, idToken string, client *mongo.Client) (userModels.LoginResponse, error) {
	userCollection := db.OpenCollection("users", client)

	claims, err := utils.ParseGoogleIDToken(idToken)

	if err != nil {
		fmt.Println("Failed to parse Google ID token:", err)
		return userModels.LoginResponse{}, utils.NewAppError("invalid Google token", 400)
	}
	fmt.Printf("Google token claims: %+v\n", claims)
	email := claims.Email
	name := claims.Name
	googleID := claims.Sub

	var user userModels.User
	err = userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return userModels.LoginResponse{}, utils.NewAppError("database error", 500)
	}

	if err == mongo.ErrNoDocuments {
		// 3️⃣ Create new Google-only user
		user = userModels.User{
			ID:              primitive.NewObjectID(),
			Name:            name,
			Email:           email,
			GoogleID:        googleID, // store Google ID
			Role:            "USER",
			IsEmailVerified: true, // Google already verified email
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		_, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			return userModels.LoginResponse{}, utils.NewAppError("failed to create user", 500)
		}
	} else {
		// 4️⃣ Link Google ID if not already linked
		if user.GoogleID == "" {
			_, err := userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
				"$set": bson.M{
					"google_id":  googleID,
					"updated_at": time.Now(),
				},
			})
			if err != nil {
				return userModels.LoginResponse{}, utils.NewAppError("failed to link Google account", 500)
			}
			user.GoogleID = googleID
		}
	}

	// 5️⃣ Generate auth tokens
	accessToken, refreshToken, _, err := utils.GenerateAuthTokens(user)
	if err != nil {
		return userModels.LoginResponse{}, utils.NewAppError("failed to generate auth tokens", 500)
	}

	return userModels.LoginResponse{
		/* User:         user.ToResponse(), */
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
