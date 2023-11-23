package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Password string             `json:"password" bson:"password"`
	Username string             `json:"username" bson:"username"`
}
