package routes

import (
	"Go-starter-template/models"
	"context"
	"net/http"
	"path/filepath"

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
