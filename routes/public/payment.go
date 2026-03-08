package public

import (
	"github.com/gin-gonic/gin"
	payment "github.com/sup25/gobuy/internal/payment"
	"go.mongodb.org/mongo-driver/mongo"
)

func PaymentRoutes(r *gin.RouterGroup, client *mongo.Client) {
	r.POST("/payments/create-intent", payment.CreatePaymentIntent(client))
}
