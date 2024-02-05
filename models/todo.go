package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID       primitive.ObjectID `json:"id" bson:"id" binding:"omitempty"`
	Title    string             `json:"title" bson:"title"`
	Priority string             `json:"priority" bson:"priority"`
	Due      time.Time          `json:"due" bson:"due"`
	Files    []string           `json:"files" bson:"files"`
	User     UserDetails
}

type UserDetails struct {
	Username string `json:"username" bson:"username"`
	Role     string `json:"role" bson:"role" binding:"oneof=superadmin admin user"`
}
