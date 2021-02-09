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
	dayArr [12]int
}

// NewDrowsinessRepository ...
func NewDrowsinessRepository(m *database.MongoDatabase, c *config.Config) *DrowsinessRepository {
	return &DrowsinessRepository{
		MONGO:  m,
		config: c,
		dayArr: [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31},
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

func (ar *DrowsinessRepository) GetNumberOfDrowsinessOnDay(startDay int, startMonth int, startYear int, endDay int, endMonth int, endYear int) (int32, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("drowsiness")
	fromHour := time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)
	toHour := time.Date(endYear, time.Month(endMonth), endDay, 23, 59, 99, 99, time.UTC)
	filter := bson.D{
		{
			"time", bson.D{
				{"$gt", fromHour},
				{"$lte", toHour},
			},
		},
	}
	cur, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	return int32(cur), nil

}

func (ar *DrowsinessRepository) GetNumberOfDrowsinessToCalendar(year int) ([]*models.DrowsinessStatCal, error) {
	year1, month, day := time.Now().Date()
	monthArr := [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	var dayArr [12]int = ar.dayArr

	if int(year)%400 == 0 || (int(year)%4 == 0 && !(int(year)%100 == 0)) {
		dayArr[1] = 29
	}
	var m int = int(month)
	if !(year == year1) {
		m = 12
	} else {
		dayArr[int(month-1)] = day
	}
	var results []*models.DrowsinessStatCal
	for j := 0; j < m; j++ {
		var result models.DrowsinessStatCal
		result.Name = monthArr[j]
		days := make([]int32, dayArr[j])
		for i := 1; i <= dayArr[j]; i++ {
			cur, err := ar.GetNumberOfDrowsinessOnDay(i, j+1, year, i, j+1, year)
			if err != nil {
				log.Fatal(err)
			}
			days[i-1] = cur
			if i == day && int(month)-1 == j {
				break
			}

		}
		result.Data = days
		results = append(results, &result)

	}

	return results, nil
}

// time
func (ar *DrowsinessRepository) GetNumberOfDrowsinessHour(day int, month int, year int, hour int32) (int32, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("drowsiness")
	fromHour := time.Date(year, time.Month(month), day, int(hour), 0, 0, 0, time.UTC)
	toHour := time.Date(year, time.Month(month), day, int(hour), 59, 99, 999, time.UTC)
	filter := bson.D{
		{
			"time", bson.D{
				{"$gt", fromHour},
				{"$lte", toHour},
			},
		},
	}
	cur, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	return int32(cur), nil

}

func (ar *DrowsinessRepository) GetNumberOfDrowsinessTimeBar(day int, month int, year int) ([]int32, error) {
	year1, month1, day1 := time.Now().Date()
	hour := 23
	var dayArr [12]int = ar.dayArr
	var mt int = 12
	var mst int = month
	days := make([]int32, hour+1)
	for y := year; y < year1+1; y++ {
		if y == year1 {
			mt = int(month1)
			dayArr[mt-1] = day1
		}
		if y%400 == 0 || (y%4 == 0 && !(y%100 == 0)) {
			dayArr[1] = 29
		} else {
			dayArr[1] = 28
		}
		for m := mst; m <= mt; m++ {
			for d := 1; d <= dayArr[m-1]; d++ {
				if day1 == d && int(month1) == m && y == year1 {
					hour = time.Now().Hour()
				}
				for i := 0; i <= hour; i++ {
					cur, err := ar.GetNumberOfDrowsinessHour(d, m, y, int32(i))
					if err != nil {
						log.Fatal(err)
					}
					days[i] = days[i] + cur
				}
			}
		}
		if y == year {
			mst = 1
		}
	}

	return days, nil
}
