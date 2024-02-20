package routes

import (
	"Go-starter-template/models"
	"context"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Parse MultipartForm
	// 10 MB limit for the entire request
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form"})
		return
	}
	username := c.PostForm("username")
	title := c.PostForm("title")
	priority := c.PostForm("priority")
	role := c.PostForm("role")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	files := form.File["files"]

	var fileNames []string
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		fileNames = append(fileNames, filename)
	}

	newTodo := models.Todo{
		ID:       primitive.NewObjectID(),
		Title:    title,
		Priority: priority,
		User: models.UserDetails{
			Username: username,
			Role:     role,
		},
		Files: fileNames,
	}

	result, err := handler.collection.InsertOne(handler.ctx, newTodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error-2": err})
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

	var todo models.Todo
	err := handler.collection.FindOne(handler.ctx, bson.M{
		"id": objectId,
	}).Decode(&todo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If todo.Files is not null or undefined, delete associated files
	if todo.Files != nil {
		for _, file := range todo.Files {
			err := os.Remove(filepath.Join("./uploads", file))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	_, err = handler.collection.DeleteOne(handler.ctx, bson.M{
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

/* a function to count all todos */
func (handler *TodosHandler) CountAllTodosHandler(c *gin.Context) {
	opts := options.Count().SetHint("_id_") //optimize search
	count, err := handler.collection.CountDocuments(handler.ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"TotalTodos": count,
	})
}

/* a function to count all todos */
func (handler *TodosHandler) CountTodosByPriorityHandler(c *gin.Context) {
	priority := c.Param("priority")
	count, err := handler.collection.CountDocuments(handler.ctx, bson.M{
		"priority": priority,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"TotalTodos": count,
	})
}

/* get todo with pagination */
func (handler *TodosHandler) ListTodosWithPagHandler(c *gin.Context) {
	// Extract query parameters for pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if invalid or not provided
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default to 10 items per page if invalid or not provided
	}

	// Calculate skip value for pagination
	skip := (page - 1) * pageSize

	// MongoDB find options for pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))

	// Perform find operation with pagination options
	cur, err := handler.collection.Find(handler.ctx, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	// Iterate over the cursor and decode todos
	todos := make([]models.Todo, 0)
	for cur.Next(handler.ctx) {
		var todo models.Todo
		cur.Decode(&todo)
		todos = append(todos, todo)
	}

	// Count total documents to calculate total pages
	totalTodos, err := handler.collection.CountDocuments(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalPages := int(math.Ceil(float64(totalTodos) / float64(pageSize)))

	// Return paginated todos along with total pages and hasMore flag
	hasMore := page < totalPages
	responseData := gin.H{
		"totalPages": totalPages,
		"data":       todos,
		"hasMore":    hasMore,
	}
	c.JSON(http.StatusOK, responseData)
}
