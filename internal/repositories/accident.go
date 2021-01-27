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

// GetHourlyAccidentOfCurrentDay ...
func (ar *AccidentRepository) GetHourlyAccidentOfCurrentDay(hour int32) ([]*models.Accident, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")

	var results []*models.Accident

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


func (ar *AccidentRepository) GetNumberOfAccidentCurrentYear() ([]*models.AccidentStatCal, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")

	year, month, day := time.Now().Date()
	monthArr := [12]string{"Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"}
	var dayArr [12]int = [12]int{31,28,31,30,31,30,31,31,30,31,30,31}

	if ( int(year)%400 == 0 || (int(year)%4==0 && !(int(year)%100==0))) {
		dayArr[1] = 29 
	}
	
	dayArr[int(month)] = day
	var results []*models.AccidentStatCal
	for j := 0; j < int(month); j++ {
		var result models.AccidentStatCal
		result.Name =  monthArr[j]
		days := make([]int32, dayArr[j])
		for i := 1; i <= dayArr[j]; i++{
			
			fromHour := time.Date(year, time.Month(j+1), i, 0, 0, 0, 0, time.UTC)
			toHour := time.Date(year, time.Month(j+1), i, 23, 59, 99, 99, time.UTC)
			filter := bson.D{
				{
					"time", bson.D{
						{"$gt", fromHour},
						{"$lte", toHour},
					},
				},
			}
			cur,err := collection.CountDocuments(context.TODO(),filter)
			if err != nil {
				log.Fatal(err)
			}
			days[i-1] = int32(cur)
			if i==day && int(month)-1 == j{
				break;
			}

		}
		result.Data =  days
		results = append(results, &result)
		
	}

	return results, nil
}