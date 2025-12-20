package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	dbName string
)

func init() {
	// Load env ONCE
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	dbName = os.Getenv("DATABASE_NAME")
	if dbName == "" {
		log.Fatal("DATABASE_NAME not set")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	Client = client
	fmt.Println("MongoDB connected")
}

func DBInstance() *mongo.Client {
	return Client
}

func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	collection := Client.Database(dbName).Collection(collectionName)
	if collection == nil {
		log.Fatalf("Failed to open collection: %s", collectionName)
	}
	return collection
}
