package services

import (
	"context"

	"github.com/sup25/gobuy/config/db"
	categoryModel "github.com/sup25/gobuy/internal/product/category/model"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProductCategoryService(
	ctx context.Context,
	req categoryModel.CategoryRequest,
	client *mongo.Client,
) (categoryModel.CategoryResponse, error) {

	categoryCollection := db.OpenCollection("categories", client)

	slug := req.Slug
	if slug == "" {
		slug = utils.Slugify(req.Name)
	}

	count, err := categoryCollection.CountDocuments(ctx, bson.M{"slug": slug})
	if err != nil {
		return categoryModel.CategoryResponse{}, utils.NewAppError("database error", 500)
	}
	if count > 0 {
		return categoryModel.CategoryResponse{}, utils.NewAppError("category with this slug already exists", 409)
	}

	category := categoryModel.Category{
		ID:       primitive.NewObjectID(),
		Name:     req.Name,
		Slug:     slug,
		IsActive: req.IsActive,
	}

	res, err := categoryCollection.InsertOne(ctx, category)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return categoryModel.CategoryResponse{}, utils.NewAppError("category already exists", 409)
		}
		return categoryModel.CategoryResponse{}, utils.NewAppError("failed to add category", 500)
	}

	category.ID = res.InsertedID.(primitive.ObjectID)

	return category.ToResponse(), nil
}
