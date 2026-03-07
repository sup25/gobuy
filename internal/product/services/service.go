package services

import (
	"context"
	"log"
	"time"

	"github.com/sup25/gobuy/config/db"
	categoryModel "github.com/sup25/gobuy/internal/product/category/model"
	productModel "github.com/sup25/gobuy/internal/product/model"
	"github.com/sup25/gobuy/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddProductService(
	ctx context.Context,
	req productModel.CreateProductRequest,
	userID string,
	userName string,
	userEmail string,
	userRole string,
	client *mongo.Client,
) (productModel.ProductResponse, error) {
	if userRole != "admin" && userRole != "manager" && userRole != "merchant" {
		return productModel.ProductResponse{}, utils.NewAppError("unauthorized: only admin, manager, or merchant can add products", 403)
	}

	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return productModel.ProductResponse{}, utils.NewAppError("invalid category ID", 400)
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return productModel.ProductResponse{}, utils.NewAppError("invalid user ID", 400)
	}

	// Fetch the full category to embed it
	categoryCollection := db.OpenCollection("categories", client)
	var category categoryModel.Category
	err = categoryCollection.FindOne(ctx, bson.M{"_id": categoryID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return productModel.ProductResponse{}, utils.NewAppError("category not found", 404)
		}
		return productModel.ProductResponse{}, utils.NewAppError("failed to fetch category", 500)
	}

	// Create embedded user info
	createdByUser := productModel.ToProductUser(userObjectID, userName, userEmail, userRole)

	product := productModel.Product{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		Images:      req.Images,
		Price:       req.Price,
		Stock:       req.Stock,
		Slug:        utils.Slugify(req.Name),
		Status:      req.Status,
		CategoryID:  categoryID,
		Category:    productModel.ToCategoryEmbed(&category),
		CreatedBy:   createdByUser,
		UpdatedBy:   createdByUser,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	productCollection := db.OpenCollection("products", client)
	res, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		return productModel.ProductResponse{}, utils.NewAppError("failed to add product", 500)
	}

	product.ID = res.InsertedID.(primitive.ObjectID)
	return product.ToResponse(), nil
}

func GetProductByIDService(
	ctx context.Context,
	productID string,
	client *mongo.Client,
) (productModel.ProductResponse, error) {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return productModel.ProductResponse{}, utils.NewAppError("invalid product ID", 400)
	}

	productCollection := db.OpenCollection("products", client)
	var product productModel.Product

	err = productCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return productModel.ProductResponse{}, utils.NewAppError("product not found", 404)
		}
		return productModel.ProductResponse{}, utils.NewAppError("failed to fetch product", 500)
	}

	return product.ToResponse(), nil
}

func GetAllProductsService(
	ctx context.Context,
	categorySlug string,
	status string,
	page string,
	limit string,
	client *mongo.Client,
) ([]productModel.ProductResponse, error) {
	productCollection := db.OpenCollection("products", client)

	filter := bson.M{}

	if categorySlug != "" {
		filter["category_slug"] = categorySlug
	}

	if status != "" {
		filter["status"] = status
	}

	log.Printf(
		"GetAllProductsService called with categorySlug=%q status=%q page=%q limit=%q filter=%+v",
		categorySlug,
		status,
		page,
		limit,
		filter,
	)

	pageInt := utils.ParseInt(page)
	limitInt := utils.ParseInt(limit)

	if pageInt <= 0 {
		pageInt = 1
	}
	if limitInt <= 0 {
		limitInt = 20
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64((pageInt - 1) * limitInt)).
		SetLimit(int64(limitInt))

	cursor, err := productCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, utils.NewAppError("failed to fetch products", 500)
	}
	defer cursor.Close(ctx)

	var products []productModel.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, utils.NewAppError("failed to decode products", 500)
	}

	var response []productModel.ProductResponse
	for _, product := range products {
		response = append(response, product.ToResponse())
	}

	return response, nil
}

func GetAllProductsBySlugService(
	ctx context.Context,
	slug string,
	client *mongo.Client,
) ([]productModel.ProductResponse, error) {
	return GetAllProductsService(ctx, slug, "", "", "", client)

}

func UpdateProductService(
	ctx context.Context,
	productID string,
	req productModel.UpdateProductRequest,
	userRole string,
	client *mongo.Client,
) (productModel.ProductResponse, error) {
	if userRole != "admin" && userRole != "manager" {
		return productModel.ProductResponse{}, utils.NewAppError("unauthorized: only admin or manager can update products", 403)
	}

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return productModel.ProductResponse{}, utils.NewAppError("invalid product ID", 400)
	}

	productCollection := db.OpenCollection("products", client)

	// Check if product exists
	var existingProduct productModel.Product
	err = productCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return productModel.ProductResponse{}, utils.NewAppError("product not found", 404)
		}
		return productModel.ProductResponse{}, utils.NewAppError("failed to fetch product", 500)
	}

	update := bson.M{
		"updated_at": time.Now(),
	}

	if req.Name != nil {
		update["name"] = *req.Name
		update["slug"] = utils.Slugify(*req.Name)
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.Images != nil {
		update["images"] = req.Images
	}
	if req.Price != nil {
		update["price"] = *req.Price
	}
	if req.Stock != nil {
		update["stock"] = *req.Stock
	}
	if req.Status != nil {
		update["status"] = *req.Status
	}
	if req.CategoryID != nil {
		categoryID, err := primitive.ObjectIDFromHex(*req.CategoryID)
		if err != nil {
			return productModel.ProductResponse{}, utils.NewAppError("invalid category ID", 400)
		}

		// Verify category exists
		categoryCollection := db.OpenCollection("categories", client)
		count, err := categoryCollection.CountDocuments(ctx, bson.M{"_id": categoryID})
		if err != nil {
			return productModel.ProductResponse{}, utils.NewAppError("failed to check category", 500)
		}
		if count == 0 {
			return productModel.ProductResponse{}, utils.NewAppError("category not found", 400)
		}

		update["category_id"] = categoryID
	}

	_, err = productCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	if err != nil {
		return productModel.ProductResponse{}, utils.NewAppError("failed to update product", 500)
	}

	// Fetch updated product
	var updatedProduct productModel.Product
	err = productCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedProduct)
	if err != nil {
		return productModel.ProductResponse{}, utils.NewAppError("failed to fetch updated product", 500)
	}

	return updatedProduct.ToResponse(), nil
}

func DeleteProductService(
	ctx context.Context,
	productID string,
	userRole string,
	client *mongo.Client,
) error {
	if userRole != "admin" && userRole != "manager" {
		return utils.NewAppError("unauthorized: only admin or manager can delete products", 403)
	}

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return utils.NewAppError("invalid product ID", 400)
	}

	productCollection := db.OpenCollection("products", client)

	result, err := productCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return utils.NewAppError("failed to delete product", 500)
	}

	if result.DeletedCount == 0 {
		return utils.NewAppError("product not found", 404)
	}

	return nil
}
