package routes

import (
	"go-starter-template/controllers"
	"go-starter-template/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// User routes
	api.POST("/register", controllers.CreateUserHandler) //http://localhost:8080/api/register
	api.POST("/login", controllers.LoginHandler)         //http://localhost:8080/api/login

	// Protected routes
	auth := api.Group("/")
	auth.Use(middleware.AuthMiddleware())

	// Todos routes
	auth.POST("/todos", controllers.CreateTodoHandler)
	auth.GET("/todos", controllers.GetTodosHandler)
	auth.PUT("/todos/:id", controllers.UpdateTodoHandler)
	auth.DELETE("/todos/:id", controllers.DeleteTodoHandler)

	//http://localhost:8080/api/admin/users
	// Admin-only routes
	admin := auth.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.GET("/users", controllers.GetAllUsersHandler)
}
