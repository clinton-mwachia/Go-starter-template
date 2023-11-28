package main

import (
	"Go-starter-template/routes"
	"context"
	"log"

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

func main() {
	router := gin.Default()

	/* users*/
	router.POST("/user/register", usersHandler.AddNewUser)
	router.GET("/users", usersHandler.ListUsersHandler)
	router.GET("user/:id", usersHandler.GetUserByIdHandler)
	router.GET("users/:role", usersHandler.ListUsersByRoleHandler)
	router.DELETE("user/:id", usersHandler.DeleteUserHandler)
	router.PUT("user/:id", usersHandler.UpdateUserHandler)
	router.PUT("/user/pwd/:id", usersHandler.UpdateUserPasswordHandler)
	router.POST("/user/login", usersHandler.SignInHandler)

	router.Run(":8080")
}
