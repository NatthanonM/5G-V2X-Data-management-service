package database

import (
	"5g-v2x-data-management-service/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Client      *mongo.Client
	databaseURI string
}

func (m *MongoDatabase) connect() {
	client, err := mongo.NewClient(options.Client().ApplyURI(m.databaseURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	m.Client = client
	fmt.Println("Connected to MongoDB!")
}

func NewMongoDatabase(c *config.Config) *MongoDatabase {
	m := &MongoDatabase{
		databaseURI: c.DatabaseURI,
	}
	m.connect()
	return m
}
