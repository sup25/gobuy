package private

import (
	"github.com/gin-gonic/gin"
	categoryController "github.com/sup25/gobuy/internal/product/category/controllers"
	productController "github.com/sup25/gobuy/internal/product/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupProductRoutes(rg *gin.RouterGroup, client *mongo.Client) {
	productRoutes := rg.Group("/products")
	{
		productRoutes.POST("", productController.AddProductController(client))
		productRoutes.PUT("/:id", productController.UpdateProductController(client))

		productRoutes.DELETE("/:id", productController.DeleteProductController(client))
	}

	categoryRoutes := rg.Group("/categories")
	{
		categoryRoutes.POST("", categoryController.AddProductCategoryController(client))

	}
}
