package repositories

import (
	"5g-v2x-data-management-service/internal/models"
	"context"
	"encoding/json"
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

func (dr *DrowsinessRepository) FindOne(filter map[string]interface{}) (*models.Drowsiness, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("drowsiness")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result *models.Drowsiness

	jsonString, err := json.Marshal(filter)
	if err != nil {
		panic(err)
	}

	var bsonFilter interface{}
	err = bson.UnmarshalExtJSON([]byte(jsonString), true, &bsonFilter)
	if err != nil {
		panic(err)
	}

	err = collection.FindOne(ctx, bsonFilter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (dr *DrowsinessRepository) Find(filter map[string]interface{}) ([]*models.Drowsiness, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("drowsiness")

	var results []*models.Drowsiness

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
