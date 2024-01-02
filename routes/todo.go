package routes

import (
	"Go-starter-template/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

/* a function to add a new todo */
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

/* a function to get all todos */
func (handler *TodosHandler) ListTodosHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	todos := make([]models.Todo, 0)
	for cur.Next(handler.ctx) {
		var todo models.Todo
		cur.Decode(&todo)
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

/* a function to get todo by id*/
func (handler *TodosHandler) GetTodoByIdHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"id": objectId,
	})
	var todo models.Todo
	err := cur.Decode(&todo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

/* a function to get todo by priority */
func (handler *TodosHandler) ListTodosByRoleHandler(c *gin.Context) {
	priority := c.Param("priority")
	cur, err := handler.collection.Find(handler.ctx, bson.M{"priority": priority})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	todos := make([]models.Todo, 0)
	for cur.Next(handler.ctx) {
		var todo models.Todo
		cur.Decode(&todo)
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

/* a function to delete a todo */
func (handler *TodosHandler) DeleteTodoHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo has been deleted"})
}

/* a function to update todo by id */
func (handler *TodosHandler) UpdateTodoHandler(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// only update title and priority
	// user role must be provided as well
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: todo.Title},
			{Key: "priority", Value: todo.Priority},
		}},
	}

	_, err = handler.collection.UpdateOne(handler.ctx, bson.M{"id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo has been updated"})
}
