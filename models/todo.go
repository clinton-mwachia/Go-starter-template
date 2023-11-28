package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID       primitive.ObjectID `json:"id" bson:"id" binding:"omitempty"`
	User     primitive.ObjectID `json:"user" bson:"user"`
	Title    string             `json:"title" bson:"title"`
	Priority string             `json:"priority" bson:"priority"`
	Due      time.Time          `json:"due" bson:"due"`
}
