package routes

import (
	"Go-starter-template/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodosHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewTodosHandler(ctx context.Context, collection *mongo.Collection) *TodosHandler {
	return &TodosHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *TodosHandler) AddNewTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newTodo = models.Todo{
		ID:       primitive.NewObjectID(),
		Title:    todo.Title,
		Priority: todo.Priority,
		User: models.UserDetails{
			Username: todo.User.Username,
			Role:     todo.User.Role,
		},
	}

	result, err := handler.collection.InsertOne(handler.ctx, newTodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{
		"message": "Todo Created",
		"ID":      result.InsertedID,
	})
}
