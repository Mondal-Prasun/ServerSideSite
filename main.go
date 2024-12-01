package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello BloodBank")

	port := ":8080"

	route := gin.Default()

	route.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"health": "health is okk",
		})
	})

	route.Static("/static", "./static")
	route.LoadHTMLGlob("templates/*")

	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Cool Cat",
		})
	})

	route.POST("/submit", func(c *gin.Context) {
		// Parse JSON payload
		var json struct {
			Counter int `json:"counter"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		fmt.Println("counter: ", json.Counter)

		// Respond with the received counter value
		c.JSON(http.StatusOK, gin.H{"counter": json.Counter})
	})

	fmt.Printf("Server is running on PORT:%s", port)

	err := route.Run(port)

	if err != nil {
		panic("Serer is not running!!!")
	}

}
