package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectDB() {
	mongoURL := os.Getenv("MONGO_URL")
	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	DB = client
	log.Println("Connected to MongoDB successfully")
}

func SetupIndexes() {
	todoCollection := DB.Database("go_starter_template").Collection("todos")
	userCollection := DB.Database("go_starter_template").Collection("users")

	// Index for Todos
	todoIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userID", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: -1}},
		},
	}

	// Index for Users
	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	if _, err := todoCollection.Indexes().CreateMany(context.Background(), todoIndexes); err != nil {
		log.Fatal("Error creating todo indexes:", err)
	}
	if _, err := userCollection.Indexes().CreateMany(context.Background(), userIndexes); err != nil {
		log.Fatal("Error creating user indexes:", err)
	}
	log.Println("MongoDB indexes created successfully")
}
