package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// port := "8080"
	// fmt.Printf("Starting server on port %s", port)

	r := gin.Default()
	r.GET("/ping", pingHandler)
	r.GET("/", basicHandler)
	r.Run(":8080")
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func basicHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}
