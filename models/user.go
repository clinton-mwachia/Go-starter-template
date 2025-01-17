package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" binding:"required"`
	Email     string             `bson:"email" binding:"required,email"`
	Password  string             `bson:"password" binding:"required"`
	Role      string             `bson:"role" binding:"required"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
