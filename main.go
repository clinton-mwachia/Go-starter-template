package main

import (
	"go-starter-template/config"
	"go-starter-template/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// http://localhost:8080/api/
func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize Gin
	r := gin.Default()

	// Apply Middleware
	r.Use(cors.Default())
	//r.Use(compression.Gzip(compression.BestCompression))

	// Connect to MongoDB
	config.ConnectDB()

	// Setup MongoDB Indexing
	config.SetupIndexes()

	// Setup Routes
	routes.SetupRoutes(r)

	// Run the server
	r.Run() // Default is :8080
}
