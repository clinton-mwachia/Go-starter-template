package helpers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthHandler holds the MongoDB collection and context
type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

var JWT_SECRET = "MYSECRET1229!!"

// AuthMiddleware is a middleware function for route authentication
func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		if tokenValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if tkn == nil || !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}
