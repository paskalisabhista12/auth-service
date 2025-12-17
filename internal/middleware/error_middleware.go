package middlewares

import (
	"auth-service/pkg/utils/exception"
	"auth-service/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*exception.AppError); ok {
				c.JSON(appErr.StatusCode, response.ErrorResponse{
					Success:   false,
					ErrorCode: appErr.Code,
					Message:   appErr.Message,
				})
				return
			}

			// Fallback: unknown error
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Success:   false,
				ErrorCode: "INTERNAL_ERROR",
				Message:   err.Error(),
			})
		}
	}
}
