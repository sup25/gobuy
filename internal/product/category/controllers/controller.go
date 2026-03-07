package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	categoryModel "github.com/sup25/gobuy/internal/product/category/model"
	"github.com/sup25/gobuy/internal/product/category/services"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProductCategoryController(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		var req categoryModel.CategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, utils.NewAppError("invalid request body", 400))
			return
		}

		res, err := services.AddProductCategoryService(ctx, req, client)
		if err != nil {
			utils.RespondWithError(c, err)
			return
		}

		utils.SendSuccess(c, http.StatusCreated, "Category created successfully", res)
	}
}
