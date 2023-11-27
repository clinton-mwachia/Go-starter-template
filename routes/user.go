package routes

import (
	"Go-starter-template/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUsersHandler(ctx context.Context, collection *mongo.Collection) *UsersHandler {
	return &UsersHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// register new user
func (handler *UsersHandler) AddNewUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()
	result, err := handler.collection.InsertOne(handler.ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{
		"message": "User Created",
		"ID":      result.InsertedID,
	})
}

// get all users
func (handler *UsersHandler) ListUsersHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	users := make([]models.User, 0)
	for cur.Next(handler.ctx) {
		var user models.User
		cur.Decode(&user)
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
