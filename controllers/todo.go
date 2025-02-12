package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go-starter-template/config"
	"go-starter-template/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create a new todo
func CreateTodoHandler(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.ID = primitive.NewObjectID()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	_, err := config.DB.Database("go_starter_template").Collection("todos").InsertOne(context.Background(), todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

// Get a todo by ID
func GetTodoByIDHandler(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var todo models.Todo
	err = config.DB.Database("go_starter_template").Collection("todos").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&todo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not Found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

// Get all todos
func GetTodosHandler(c *gin.Context) {
	cursor, err := config.DB.Database("go_starter_template").Collection("todos").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var todos []models.Todo
	for cursor.Next(context.Background()) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse todos"})
			return
		}
		todos = append(todos, todo)
	}
	c.JSON(http.StatusOK, todos)
}

// Get all todos by user ID
func GetTodosByUserIDHandler(c *gin.Context) {
	userID := c.Param("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	cursor, err := config.DB.Database("go_starter_template").Collection("todos").Find(context.Background(), bson.M{"userID": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var todos []models.Todo
	for cursor.Next(context.Background()) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse todos"})
			return
		}
		todos = append(todos, todo)
	}
	c.JSON(http.StatusOK, todos)
}

// Get paginated todos
func GetTodosPaginatedHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)
	findOptions := options.Find().SetSkip(skip).SetLimit(int64(limit))

	userID := c.Param("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	cursor, err := config.DB.Database("go_starter_template").Collection("todos").Find(context.Background(), bson.M{"userID": objID}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var todos []models.Todo
	for cursor.Next(context.Background()) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse todos"})
			return
		}
		todos = append(todos, todo)
	}

	// Count the total todos for the user
	count, err := config.DB.Database("go_starter_template").Collection("todos").CountDocuments(context.Background(), bson.M{"userID": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count todos"})
		return
	}

	// Respond with paginated todos and metadata
	c.JSON(http.StatusOK, gin.H{
		"todos":      todos,
		"page":       page,
		"limit":      limit,
		"total":      count,
		"totalPages": (count + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

// Count todos by user
func CountTodosByUserHandler(c *gin.Context) {
	userID := c.Param("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	count, err := config.DB.Database("go_starter_template").Collection("todos").CountDocuments(context.Background(), bson.M{"userID": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID":     userID,
		"totalTodos": count,
	})
}

// Update a todo
func UpdateTodoHandler(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["updatedAt"] = time.Now()

	_, err = config.DB.Database("go_starter_template").Collection("todos").UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": updates})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

// Delete a todo
func DeleteTodoHandler(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = config.DB.Database("go_starter_template").Collection("todos").DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
