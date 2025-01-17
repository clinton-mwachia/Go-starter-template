package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title" binding:"required"`
	Completed bool               `bson:"completed"`
	UserID    primitive.ObjectID `bson:"userID" binding:"required"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
