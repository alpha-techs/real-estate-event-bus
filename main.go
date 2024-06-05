package main

import (
	"github.com/gin-gonic/gin"
	"real-estate-bot/service"
)

type UpdateBeRequest struct {
	Version string `json:"version"`
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/update-be", func(c *gin.Context) {
		var json UpdateBeRequest

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Call the UpdateBe function from the service package
		err := service.UpdateBe(
			json.Version,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"version": json.Version,
		})
	})

	err := r.Run(":9100")
	if err != nil {
		return
	}
}
