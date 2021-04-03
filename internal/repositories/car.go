package repositories

import (
	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/infrastructures/database"
	"5g-v2x-data-management-service/internal/models"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// CarRepository ...
type CarRepository struct {
	MONGO  *database.MongoDatabase
	config *config.Config
}

// NewCarRepository ...
func NewCarRepository(m *database.MongoDatabase, c *config.Config) *CarRepository {
	return &CarRepository{
		MONGO:  m,
		config: c,
	}
}

func (cr *CarRepository) Create(car *models.Car) (*string, error) {
	collection := cr.MONGO.Client.Database(cr.config.DatabaseName).Collection("car")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	id := uuid.New().String()
	car.CarID = id
	car.RegisteredAt = time.Now().UTC()
	defer cancel()
	_, err := collection.InsertOne(ctx, car)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (cr *CarRepository) FindOne(carID, vehRegNo *string) (*models.Car, error) {
	collection := cr.MONGO.Client.Database(cr.config.DatabaseName).Collection("car")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result *models.Car

	filterDeleted := bson.M{
		"deleted_at": bson.M{
			"$eq": nil,
		},
	}

	inputFilter := bson.M{}

	if carID != nil {
		inputFilter = bson.M{
			"_id": *carID,
		}
	}
	if vehRegNo != nil {
		inputFilter = bson.M{
			"vehicle_registration_number": *vehRegNo,
		}
	}

	filter := bson.M{
		"$and": []bson.M{
			filterDeleted,
			inputFilter,
		},
	}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cr *CarRepository) FindAll() ([]*models.Car, error) {
	collection := cr.MONGO.Client.Database(cr.config.DatabaseName).Collection("car")

	var results []*models.Car

	filter := bson.D{{}}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.Car
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

func (cr *CarRepository) Update(updateCar *models.Car) error {
	collection := cr.MONGO.Client.Database(cr.config.DatabaseName).Collection("car")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bsonFilter := bson.M{"_id": updateCar.CarID}
	updateInput := bson.D{}
	if updateCar.CarDetail != nil {
		updateInput = append(updateInput, bson.E{"car_detail", *updateCar.CarDetail})
	}
	if updateCar.VehicleRegistrationNumber != nil {
		updateInput = append(updateInput, bson.E{"vehicle_registration_number", *updateCar.VehicleRegistrationNumber})
	}

	bsonUpdate := bson.D{
		{
			"$set", updateInput,
		},
	}

	_, err := collection.UpdateOne(ctx, bsonFilter, bsonUpdate)

	if err != nil {
		return err
	}
	return nil
}
