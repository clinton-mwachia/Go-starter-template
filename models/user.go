package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"id" binding:"omitempty"`
	Password string             `json:"password" bson:"password" binding:"required"`
	Username string             `json:"username" bson:"username" binding:"required"`
	Role     string             `json:"role" bson:"role" default:"user" binding:"required,oneof=superadmin admin user"`
}
