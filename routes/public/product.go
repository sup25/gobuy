package public

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sup25/gobuy/internal/product/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProductPublicRoutes(r *gin.RouterGroup, client *mongo.Client) {
	r.GET("/products", controller.GetAllProductsController(client))
	r.GET("/products/:slug", controller.GetAllProductsBySlugController(client))

}
