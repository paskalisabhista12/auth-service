package response

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success   bool `json:"success"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// Build success response
func Success(c *gin.Context, statusCode int, data interface{}, message ...string) {
    res := SuccessResponse{
        Success: true,
        Data:    data,
    }
    if len(message) > 0 {
        res.Message = message[0]
    }
    c.JSON(statusCode, res)
}

func Error(c *gin.Context, statusCode int, errorCode, message string) {
    c.JSON(statusCode, ErrorResponse{
        Success:   false,
        ErrorCode: errorCode,
        Message:   message,
    })
}
