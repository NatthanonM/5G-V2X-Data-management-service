package repositories

import (
	"5g-v2x-data-management-service/internal/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/infrastructures/database"
)

// AccidentRepository ...
type AccidentRepository struct {
	MONGO  *database.MongoDatabase
	config *config.Config
}

// NewAccidentRepository ...
func NewAccidentRepository(m *database.MongoDatabase, c *config.Config) *AccidentRepository {
	return &AccidentRepository{
		MONGO:  m,
		config: c,
	}
}

// FindAll ...
func (ar *AccidentRepository) FindAll() ([]*models.Accident, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")

	var results []*models.Accident

	filter := bson.D{{}}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.Accident
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return results, nil
}
