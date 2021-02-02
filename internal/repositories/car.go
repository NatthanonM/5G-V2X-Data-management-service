package repositories

import (
	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/infrastructures/database"
	"5g-v2x-data-management-service/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
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
