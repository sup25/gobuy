package payment

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePaymentIntent(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var body struct {
			ProductID string `json:"product_id"`
			Amount    int64  `json:"amount"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(body.Amount),
			Currency: stripe.String("usd"),
			Metadata: map[string]string{
				"product_id": body.ProductID,
			},
		}

		pi, err := paymentintent.New(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"client_secret": pi.ClientSecret,
		})
	}
}
