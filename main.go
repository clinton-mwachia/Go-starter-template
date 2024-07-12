package main

import (
	"Go-starter-template/helpers"
	"Go-starter-template/routes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var client *mongo.Client
var err error
var MONGO_URI = "mongodb://127.0.0.1:27017/todo"
var usersHandler *routes.UsersHandler
var todosHandler *routes.TodosHandler
var authHandler *helpers.AuthHandler

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(MONGO_URI))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")

	users_collection := client.Database("todo").Collection("users")
	usersHandler = routes.NewUsersHandler(ctx, users_collection)

	todos_collection := client.Database("todo").Collection("todos")
	todosHandler = routes.NewTodosHandler(ctx, todos_collection)
}

func IndeHandler(c *gin.Context) {
	c.File("index.html")
}

func main() {
	// Set up log file
	logFilePath := "logs/app.log"
	err := ensureDirectory(filepath.Dir(logFilePath))
	if err != nil {
		fmt.Println("Error ensuring directory:", err)
		return
	}

	f, err := os.Create(logFilePath)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer f.Close()

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Set up log to write to both console and file
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()

	/* server stativ files*/
	router.Static("/public", "./public")
	/* index router */
	router.GET("/", IndeHandler)

	/* users*/
	v1 := router.Group("/users")
	{
		v1.POST("/register", usersHandler.AddNewUser)
		v1.GET("/", usersHandler.ListUsersHandler)
		v1.GET("/:id", usersHandler.GetUserByIdHandler)
		v1.GET("/role/:role", usersHandler.ListUsersByRoleHandler)
		v1.DELETE("/:id", usersHandler.DeleteUserHandler)
		v1.PUT("/:id", usersHandler.UpdateUserHandler)
		v1.PUT("/pwd/:id", usersHandler.UpdateUserPasswordHandler)
		v1.POST("/login", usersHandler.SignInHandler)
	}
	/* users */

	/* todos */
	v2 := router.Group("/todos", authHandler.AuthMiddleware())
	{
		v2.POST("/register", todosHandler.AddNewTodo)
		v2.GET("/", todosHandler.ListTodosHandler)
		v2.GET("/:id", todosHandler.GetTodoByIdHandler)
		v2.GET("/priority/:priority", todosHandler.ListTodosByRoleHandler)
		v2.DELETE("/:id", todosHandler.DeleteTodoHandler)
		v2.PUT("/:id", todosHandler.UpdateTodoHandler)
		v2.GET("/count", todosHandler.CountAllTodosHandler)
		v2.GET("/count/:priority", todosHandler.CountTodosByPriorityHandler)
		v2.GET("/paginator", todosHandler.ListTodosWithPagHandler)
	}

	router.Run(":8080")
}

// ensureDirectory creates the directory if it doesn't exist.
func ensureDirectory(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
