package routes

import (
	"Go-starter-template/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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
