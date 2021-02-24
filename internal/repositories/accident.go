package repositories

import (
	"5g-v2x-data-management-service/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/infrastructures/database"
)

// AccidentRepository ...
type AccidentRepository struct {
	MONGO  *database.MongoDatabase
	config *config.Config
	dayArr [12]int
}

// NewAccidentRepository ...
func NewAccidentRepository(m *database.MongoDatabase, c *config.Config) *AccidentRepository {
	return &AccidentRepository{
		MONGO:  m,
		config: c,
		dayArr: [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31},
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

func (ar *AccidentRepository) Find(filter primitive.D) ([]*models.Accident, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")

	var results []*models.Accident

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

func (ar *AccidentRepository) GetNumberOfAccidentHour(day int, month int, year int, hour int32) (int32, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")
	fromHour := time.Date(year, time.Month(month), day, int(hour), 0, 0, 0, time.UTC)
	toHour := time.Date(year, time.Month(month), day, int(hour)+1, 0, 0, 0, time.UTC)
	filter := bson.D{
		{
			"time", bson.D{
				{"$gte", fromHour},
				{"$lt", toHour},
			},
		},
	}
	cur, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	return int32(cur), nil

}

func (ar *AccidentRepository) GetNumberOfAccidentTimeBar(day int, month int, year int) ([]int32, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")
	t := time.Now()
	h := t.Hour() + 1
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		h = 24
	}
	thTimeZone, _ := time.LoadLocation("Asia/Bangkok")
	fromTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, thTimeZone).UTC()
	toTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 999, thTimeZone).UTC()

	filter := bson.D{{
		"time", bson.D{
			{"$gte", fromTime},
			{"$lt", toTime},
		},
	}}
	matchStage := bson.D{{"$match", filter}}

	projectStage := bson.D{{
		"$project", bson.M{
			"h": bson.M{"$hour": "$time"},
		},
	}}

	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{
			{"hour", "$h"},
		}},
		{"total", bson.D{{"$sum", 1}}},
	}}}

	cur, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, groupStage})
	if err != nil {
		log.Fatal(err)
	}
	countEachHour := make([]int32, h)
	for cur.Next(context.TODO()) {
		var elem models.NumberOfAccident
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		countEachHour[(int(elem.ID.Hour)+7)%24] = elem.Total
	}
	return countEachHour, nil
}

func (ar *AccidentRepository) GetNumberOfAccidentDay(startDay int, startMonth int, startYear int, endDay int, endMonth int, endYear int) (int32, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")
	fromHour := time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)
	toHour := time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.UTC)
	filter := bson.D{
		{
			"time", bson.D{
				{"$gte", fromHour},
				{"$lt", toHour},
			},
		},
	}
	cur, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	return int32(cur), nil

}

func (ar *AccidentRepository) GetNumberOfAccidentToCalendar(year int) ([]*models.AccidentStatCal, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")
	t := time.Now()
	thTimeZone, _ := time.LoadLocation("Asia/Bangkok")
	fromTime := time.Date(year, time.Month(0), 1, 0, 0, 0, 0, thTimeZone).UTC()
	toTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, thTimeZone).UTC()
	year1 := toTime.Year()
	monthArr := [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	var dayArr [12]int = ar.dayArr
	if int(year)%400 == 0 || (int(year)%4 == 0 && !(int(year)%100 == 0)) {
		dayArr[1] = 29
	}
	var m int = int(t.Month())
	if !(year == year1) {
		m = 12
		toTime = time.Date(year, time.Month(11), 31, 23, 59, 99, 999, thTimeZone).UTC()
	} else {
		dayArr[m-1] = t.Day()
	}
	var results []*models.AccidentStatCal
	for j := 0; j < m; j++ {
		var result models.AccidentStatCal
		result.Name = monthArr[j]
		result.Data = make([]int32, dayArr[j])
		results = append(results, &result)
	}

	var layoutISO string = "2006-01-02"
	filter := bson.D{{
		"time", bson.D{
			{"$gte", fromTime},
			{"$lt", toTime},
		},
	}}

	matchStage := bson.D{{"$match", filter}}

	projectStage := bson.D{{
		"$project", bson.M{
			"d": bson.M{
				"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$time", "timezone": "Asia/Bangkok"},
			},
		},
	}}

	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{
			{"date", "$d"},
		}},
		{"total", bson.D{{"$sum", 1}}},
	}}}

	cur, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, groupStage})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.NumberOfAccidentDate
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		t, _ := time.Parse(layoutISO, string(elem.ID.Date))
		results[int(t.Month())-1].Data[t.Day()-1] = elem.Total
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return results, nil
}

func (ar *AccidentRepository) GetNumberOfAccidentStreet(startDay int, startMonth int, startYear int) ([]*models.NumberOfAccidentRoad, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")
	year, month, day := time.Now().Date()
	thTimeZone, _ := time.LoadLocation("Asia/Bangkok")
	fromHour := time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, thTimeZone).UTC()
	toHour := time.Date(year, month, day, 23, 59, 99, 999, thTimeZone).UTC()
	var m []*models.NumberOfAccidentRoad

	groupStage := bson.D{
		{"$group", bson.D{{"_id", "$road"}, {"total", bson.D{{"$sum", 1}}}}},
	}
	sortStage := bson.D{
		{"$sort", bson.D{{"total", -1}}},
	}

	matchStage := bson.D{{"$match", bson.D{{
		"time", bson.D{
			{"$gte", fromHour},
			{"$lt", toHour},
		},
	}}}}

	cur, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage,sortStage})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem models.NumberOfAccidentRoad
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		
		m = append(m, &elem)
		// m["N/A"] = elem.Total
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())
	fmt.Println(m)

	return m, nil
}

func (ar *AccidentRepository) GetAccidentStatGroupByHour(from, to *timestamppb.Timestamp, driverUsername *string) ([24]int32, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("accident")
	fromTime := time.Date(1970, time.Month(0), 0, 0, 0, 0, 0, time.UTC)
	toTime := time.Now()
	if from != nil {
		fromTime = from.AsTime()
	}
	if to != nil {
		toTime = to.AsTime()
	}

	filter := bson.D{{
		"time", bson.D{
			{"$gte", fromTime},
			{"$lt", toTime},
		},
	}}
	if driverUsername != nil {
		filter = append(filter, bson.E{
			"username", *driverUsername,
		})
	}
	matchStage := bson.D{{"$match", filter}}

	projectStage := bson.D{{
		"$project", bson.M{
			"h": bson.M{"$hour": bson.D{{"date", "$time"}, {"timezone", "Asia/Bangkok"}}},
		},
	}}

	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{
			{"hour", "$h"},
		}},
		{"total", bson.D{{"$sum", 1}}},
	}}}

	cur, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, groupStage})
	if err != nil {
		log.Fatal(err)
	}

	countEachHour := [24]int32{}
	for cur.Next(context.TODO()) {
		var elem models.NumberOfAccident
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		countEachHour[(int(elem.ID.Hour))%24] = elem.Total
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return countEachHour, nil
}
