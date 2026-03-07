package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	productModels "github.com/sup25/gobuy/internal/product/model"
	"github.com/sup25/gobuy/internal/product/services"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddProductController creates a new product
func AddProductController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user info from JWT claims
		userID := c.GetString("user_id")
		userName := c.GetString("name")
		userEmail := c.GetString("email")
		userRole := c.GetString("role")

		fmt.Printf("UserID: %s, UserName: %s, UserEmail: %s, UserRole: %s\n",
			userID, userName, userEmail, userRole)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var req productModels.CreateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError(err.Error(), 400))
			return
		}

		res, err := services.AddProductService(ctx, req, userID, userName, userEmail, userRole, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusCreated, "Product created successfully", res)
	}
}

// GetProductByIDController retrieves a single product by ID
func GetProductByIDController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		res, err := services.GetProductByIDService(ctx, productID, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusOK, "Product retrieved successfully", res)
	}
}

// GetAllProductsController retrieves all products with optional filters
func GetAllProductsController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		categorySlug := c.Param("slug")
		status := c.Query("status")
		page := c.Query("page")
		limit := c.Query("limit")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		res, err := services.GetAllProductsService(ctx, categorySlug, status, page, limit, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusOK, "Products retrieved successfully", res)
	}
}

func GetAllProductsBySlugController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		res, err := services.GetAllProductsBySlugService(ctx, slug, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusOK, "Products retrieved successfully", res)
	}

}

// UpdateProductController updates an existing product
func UpdateProductController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		productID := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var req productModels.UpdateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError(err.Error(), 400))
			return
		}

		res, err := services.UpdateProductService(ctx, productID, req, userRole, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusOK, "Product updated successfully", res)
	}
}

// DeleteProductController deletes a product
func DeleteProductController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		productID := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := services.DeleteProductService(ctx, productID, userRole, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusOK, "Product deleted successfully", nil)
	}
}
