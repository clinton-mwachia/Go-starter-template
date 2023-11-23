package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Password string             `json:"password" bson:"password" validate:"required"`
	Username string             `json:"username" bson:"username" validate:"required"`
	Role     string             `json:"role" bson:"role" default:"user" validate:"required,oneof=superadmin admin user"`
}
