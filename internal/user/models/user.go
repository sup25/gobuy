package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`
	Email           string             `bson:"email" json:"email"`
	Password        string             `bson:"password" json:"-"`
	Role            string             `bson:"role" json:"role"` // USER | MANAGER | ADMIN
	Permissions     []string           `bson:"permissions" json:"permissions"`
	IsEmailVerified bool               `bson:"is_email_verified" json:"is_email_verified"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
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
