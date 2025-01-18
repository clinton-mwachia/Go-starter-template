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
	api.GET("/todos", controllers.GetTodosHandler)
	api.PUT("/todo/:id", controllers.UpdateTodoHandler)
	api.DELETE("/todo/:id", controllers.DeleteTodoHandler)
	api.POST("/todo", controllers.CreateTodoHandler)
	api.GET("/todos/count/:userID", controllers.CountTodosByUserHandler)

	// Users routes
	api.GET("/users", controllers.GetAllUsersHandler)
	api.PUT("/user/:id", controllers.UpdateUserHandler)
	api.DELETE("/user/:id", controllers.DeleteUserHandler)
	api.POST("/user", controllers.CreateUserHandler)
	api.GET("/user/count", controllers.CountUsersHandler)

	//http://localhost:8080/api/admin/users
	// Admin-only routes
	admin := auth.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.GET("/users", controllers.GetAllUsersHandler)
}
