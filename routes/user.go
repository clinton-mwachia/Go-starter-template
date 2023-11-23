package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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

func (handler *UsersHandler) GetUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello from users",
	})
}
