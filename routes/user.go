package routes

import (
	"Go-starter-template/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

var secret = "mysecret111"

func NewUsersHandler(ctx context.Context, collection *mongo.Collection) *UsersHandler {
	return &UsersHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// register new user
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

// get all users
func (handler *UsersHandler) ListUsersHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	users := make([]models.User, 0)
	for cur.Next(handler.ctx) {
		var user models.User
		cur.Decode(&user)
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// get user by id
func (handler *UsersHandler) GetUserByIdHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"id": objectId,
	})
	var user models.User
	err := cur.Decode(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// get users by role
func (handler *UsersHandler) ListUsersByRoleHandler(c *gin.Context) {
	role := c.Param("role")
	cur, err := handler.collection.Find(handler.ctx, bson.M{"role": role})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	users := make([]models.User, 0)
	for cur.Next(handler.ctx) {
		var user models.User
		cur.Decode(&user)
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// delete a user
func (handler *UsersHandler) DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User has been deleted"})
}

// update user details
func (handler *UsersHandler) UpdateUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "username", Value: user.Username},
			{Key: "role", Value: user.Role},
		}},
	}

	_, err = handler.collection.UpdateOne(handler.ctx, bson.M{"id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User has been updated"})
}

// update user password
func (handler *UsersHandler) UpdateUserPasswordHandler(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password", Value: string(hashedPassword)},
		}},
	}

	_, err = handler.collection.UpdateOne(handler.ctx, bson.M{"id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User password updated"})
}

// log in user
func (handler *UsersHandler) SignInHandler(c *gin.Context) {
	var userInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind the JSON request to the userInput variable
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": userInput.Username,
	})

	err := cur.Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.ID.Hex(),
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}).SignedString([]byte(secret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token,
		"expiry": time.Now().Add(time.Hour * 1).Unix()})
}
