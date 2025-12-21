package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                     string             `bson:"name" json:"name"`
	Email                    string             `bson:"email" json:"email"`
	Password                 string             `bson:"password" json:"-"`
	GoogleID                 string             `bson:"google_id,omitempty" json:"google_id,omitempty"`
	Role                     string             `bson:"role" json:"role"`
	Permissions              []string           `bson:"permissions" json:"permissions"`
	IsEmailVerified          bool               `bson:"is_email_verified" json:"is_email_verified"`
	CreatedAt                time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt                time.Time          `bson:"updated_at" json:"updated_at"`
	PasswordResetToken       string             `bson:"password_reset_token,omitempty" json:"-"`
	PasswordResetExpires     time.Time          `bson:"password_reset_expires,omitempty" json:"-"`
	EmailVerificationToken   string             `bson:"email_verification_token,omitempty" json:"-"`
	EmailVerificationExpires time.Time          `bson:"email_verification_expires,omitempty" json:"-"`
}

type LoginResponse struct {
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

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Sub           string `json:"sub"`
	jwt.RegisteredClaims
}

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
