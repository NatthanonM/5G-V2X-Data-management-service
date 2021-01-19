package repositories

import (
	"context"
	"fmt"
	"time"

	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/infrastructures/database"
)

// CRUDRepository ...
type CRUDRepository struct {
	MONGO  *database.MongoDatabase
	config *config.Config
}

// NewCRUDRepository ...
func NewCRUDRepository(m *database.MongoDatabase, c *config.Config) *CRUDRepository {
	return &CRUDRepository{
		MONGO:  m,
		config: c,
	}
}

// Create ...
func (crud *CRUDRepository) Create(collectionName string, m interface{}) (string, error) {
	collection := crud.MONGO.Client.Database(crud.config.DatabaseName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, m)
	if err != nil {
		return "", err
	}
	id := fmt.Sprintf("%s", result.InsertedID)
	return id, nil
}
