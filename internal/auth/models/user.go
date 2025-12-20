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
