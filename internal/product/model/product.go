package model

import (
	"time"

	categoryModel "github.com/sup25/gobuy/internal/product/category/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Embedded lightweight category for product listings
type ProductCategory struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
	Slug string             `bson:"slug" json:"slug"`
}

// Embedded user info for tracking who created/updated the product
type ProductUser struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Role  string             `bson:"role" json:"role"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	Images      []string           `bson:"images,omitempty" json:"images"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	Slug        string             `bson:"slug" json:"slug"`
	Status      string             `bson:"status" json:"status"`
	CategoryID  primitive.ObjectID `bson:"category_id" json:"category_id"`
	Category    *ProductCategory   `bson:"category,omitempty" json:"category,omitempty"`
	MerchantID  primitive.ObjectID `bson:"merchant_id" json:"merchant_id"`
	CreatedBy   *ProductUser       `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy   *ProductUser       `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Stock       int      `json:"stock" binding:"gte=0"`
	CategoryID  string   `json:"category_id" binding:"required"`
	Status      string   `json:"status"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Images      []string `json:"images,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
	Status      *string  `json:"status,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
}

type ProductResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Images      []string           `json:"images,omitempty"`
	Price       float64            `json:"price"`
	Stock       int                `json:"stock"`
	Slug        string             `json:"slug"`
	Status      string             `json:"status"`
	CategoryID  primitive.ObjectID `json:"category_id"`
	MerchantID  primitive.ObjectID `json:"merchant_id"`
	Category    *ProductCategory   `json:"category,omitempty"`
	Merchant    *MerchantInfo      `json:"merchant,omitempty"` // Optional
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type MerchantInfo struct {
	ID    primitive.ObjectID `json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
}

func (p *Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Images:      p.Images,
		Price:       p.Price,
		Stock:       p.Stock,
		Slug:        p.Slug,
		Status:      p.Status,
		CategoryID:  p.CategoryID,
		MerchantID:  p.MerchantID,
		Category:    p.Category,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// Helper to convert full Category to ProductCategory
func ToCategoryEmbed(cat *categoryModel.Category) *ProductCategory {
	if cat == nil {
		return nil
	}
	return &ProductCategory{
		ID:   cat.ID,
		Name: cat.Name,
		Slug: cat.Slug,
	}
}

// Helper to create ProductUser from user details
func ToProductUser(userID primitive.ObjectID, name, email, role string) *ProductUser {
	return &ProductUser{
		ID:    userID,
		Name:  name,
		Email: email,
		Role:  role,
	}
}
