package main

import (
	"orderApp/configs"
	"orderApp/order"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config.Connect()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the order App",
		})
	})

	order.OrderRoutes(r)

	r.Run()
}
