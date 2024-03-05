package main

import (
	"Go-starter-template/routes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
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

	// Set up log to write to both console and file
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()

	/* server stativ files*/
	router.Static("/public", "./public")
	/* index router */
	router.GET("/", IndeHandler)

	/* users*/
	router.POST("/user/register", usersHandler.AddNewUser)
	router.GET("/users", usersHandler.ListUsersHandler)
	router.GET("user/:id", usersHandler.GetUserByIdHandler)
	router.GET("users/:role", usersHandler.ListUsersByRoleHandler)
	router.DELETE("user/:id", usersHandler.DeleteUserHandler)
	router.PUT("user/:id", usersHandler.UpdateUserHandler)
	router.PUT("/user/pwd/:id", usersHandler.UpdateUserPasswordHandler)
	router.POST("/user/login", usersHandler.SignInHandler)
	/* users */

	/* todos */
	router.POST("/todo/register", todosHandler.AddNewTodo)
	router.GET("/todos", todosHandler.ListTodosHandler)
	router.GET("/todo/:id", todosHandler.GetTodoByIdHandler)
	router.GET("/todos/:priority", todosHandler.ListTodosByRoleHandler)
	router.DELETE("/todo/:id", todosHandler.DeleteTodoHandler)
	router.PUT("/todo/:id", todosHandler.UpdateTodoHandler)
	router.GET("/todos/count", todosHandler.CountAllTodosHandler)
	router.GET("/todos/count/:priority", todosHandler.CountTodosByPriorityHandler)
	router.GET("/todos/paginator", todosHandler.ListTodosWithPagHandler)

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
