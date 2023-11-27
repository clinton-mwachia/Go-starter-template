package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"id" binding:"omitempty"`
	Password string             `json:"password" bson:"password"`
	Username string             `json:"username" bson:"username"`
	Role     string             `json:"role" bson:"role" binding:"oneof=superadmin admin user"`
}
