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
	auth.GET("/todos", controllers.GetTodosHandler)
	auth.PUT("/todo/:id", controllers.UpdateTodoHandler)
	auth.DELETE("/todo/:id", controllers.DeleteTodoHandler)
	auth.GET("/todo/:id", controllers.GetTodoByIDHandler)
	auth.POST("/todo", controllers.CreateTodoHandler)
	auth.GET("/todos/count/:userID", controllers.CountTodosByUserHandler)

	// Users routes
	auth.PUT("/user/:id", controllers.UpdateUserHandler)
	auth.DELETE("/user/:id", controllers.DeleteUserHandler)
	auth.POST("/user", controllers.CreateUserHandler)
	auth.GET("/user/count", controllers.CountUsersHandler)
	auth.GET("/user/:id", controllers.GetUserByIDHandler)

	//http://localhost:8080/api/admin/users
	// Admin-only routes
	admin := auth.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.GET("/users", controllers.GetAllUsersHandler)
}
