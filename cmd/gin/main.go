package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
}

func main() {
	// Create Gin router
	router := gin.Default()

	// Register Routes
	router.GET("/", homePage)

	// Start the server
	router.Run()
}
