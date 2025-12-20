package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ Request DTOs (Data Transfer Objects) ============

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

// ============ Response DTOs ============

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type UserResponse struct {
	ID              primitive.ObjectID `json:"id"`
	Name            string             `json:"name"`
	Email           string             `json:"email"`
	Role            string             `json:"role"`
	Permissions     []string           `json:"permissions"`
	IsEmailVerified bool               `json:"is_email_verified"`
	CreatedAt       time.Time          `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"` // Optional for refresh endpoint
}

type MessageResponse struct {
	Message string `json:"message"`
}

// ============ Database Models ============

type RefreshToken struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	TokenHash  string             `bson:"token_hash" json:"-"`
	ExpiresAt  time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	LastUsedAt time.Time          `bson:"last_used_at" json:"last_used_at"`
	IsRevoked  bool               `bson:"is_revoked" json:"is_revoked"`
	DeviceInfo *string            `bson:"device_info,omitempty" json:"device_info,omitempty"`
	IPAddress  *string            `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
}

type PasswordResetToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	TokenHash string             `bson:"token_hash" json:"-"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	IsUsed    bool               `bson:"is_used" json:"is_used"`
}

// ============ Helper Methods ============

// Convert User to UserResponse (excludes sensitive fields)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Role:            u.Role,
		Permissions:     u.Permissions,
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt,
	}
}
