package main

import (
	"github.com/gin-gonic/gin"
	"auth-service/infra/db"
)

func main() {

	db.Connect()

	r := gin.Default()

	// Simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Run server on port 8080
	r.Run(":8080")
}
