package repositories

import (
	"5g-v2x-data-management-service/internal/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/infrastructures/database"
)

// DrowsinessRepository ...
type DrowsinessRepository struct {
	MONGO  *database.MongoDatabase
	config *config.Config
}

// NewDrowsinessRepository ...
func NewDrowsinessRepository(m *database.MongoDatabase, c *config.Config) *DrowsinessRepository {
	return &DrowsinessRepository{
		MONGO:  m,
		config: c,
	}
}

// GetHourlyDrowsinessOfCurrentDay ...
func (dr *DrowsinessRepository) GetHourlyDrowsinessOfCurrentDay(hour int32) ([]*models.Drowsiness, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("drowsiness")

	var results []*models.Drowsiness

	year, month, day := time.Now().Date()
	fromHour := time.Date(year, month, day, int(hour), 0, 0, 0, time.UTC)
	toHour := time.Date(year, month, day, int(hour+1), 0, 0, 0, time.UTC)

	filter := bson.D{
		{
			"time", bson.D{
				{"$gt", fromHour},
				{"$lte", toHour},
			},
		},
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.Drowsiness
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
