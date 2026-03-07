// routes/public/public.go
package public

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func PublicRoute(router *gin.Engine, client *mongo.Client) {
	public := router.Group("/api/v1")
	AuthPublicRoutes(public, client)
	ProductPublicRoutes(public, client)
}
