package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize Gin router
	r := gin.Default()

	r.Run(":8080")
}
