package middlewares

import (
	"auth-service/pkg/utils"
	"context"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const TrxIDKey ctxKey = "transaction_id"

func TransactionIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		trxID := c.GetHeader("X-Transaction-ID")
		if trxID == "" {
			trxID = utils.GenerateTransactionID()
		}

		ctx := context.WithValue(c.Request.Context(), TrxIDKey, trxID)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Transaction-ID", trxID)
		c.Next()
	}
}