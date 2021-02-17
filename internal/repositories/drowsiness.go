package repositories

import (
	"5g-v2x-data-management-service/internal/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"

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

func (dr *DrowsinessRepository) Find(filter primitive.D) ([]*models.Drowsiness, error) {
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

func (ar *DrowsinessRepository) GetNumberOfDrowsinessTimeBar(day int, month int, year int) ([]int32, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("drowsiness")
	t := time.Now()
	h := t.Hour()+1
	if(t.Year() != year ||int(t.Month()) != month||t.Day() != day){
		h = 24
	}
	thTimeZone, _ := time.LoadLocation("Asia/Bangkok")
	fromTime := time.Date(year,time.Month(month), day, 0, 0, 0, 0, thTimeZone).UTC()
	toTime := time.Date(t.Year(), t.Month(), t.Day(),t.Hour(), 59, 59, 999, thTimeZone).UTC()
	
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
		var elem models.NumberOfDrowsiness
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		countEachHour[(int(elem.ID.Hour)+7)%24] = elem.Total
	}
	return countEachHour, nil
}

// time
func (dr *DrowsinessRepository) GetNumberOfDrowsinessHour(day int, month int, year int, hour int32) (int32, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("drowsiness")
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

func (ar *DrowsinessRepository) GetNumberOfDrowsinessToCalendar(year int) ([]*models.DrowsinessStatCal, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("drowsiness")
	t := time.Now()
	thTimeZone, _ := time.LoadLocation("Asia/Bangkok")
	fromTime := time.Date(year,time.Month(0), 1, 0, 0, 0, 0, thTimeZone).UTC()
	toTime := time.Date(t.Year(), t.Month(), t.Day(),23, 59, 59, 999, thTimeZone).UTC()
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
	var results []*models.DrowsinessStatCal
	for j := 0; j < m; j++ {
		var result models.DrowsinessStatCal
		result.Name = monthArr[j]
		result.Data = make([]int32, dayArr[j])
		results = append(results, &result)
	}
	
	var layoutISO string= "2006-01-02"
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
				"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$time",    "timezone": "Asia/Bangkok",},
			},
		},
	}}

	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{
			{"date", "$d"},		
		}},
		{"total", bson.D{{"$sum", 1}}},
	}}}
	
	cur, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage,projectStage,groupStage})
	if err != nil {
		log.Fatal(err)
	}
	
	for cur.Next(context.TODO()) {
		var elem models.NumberOfDrowsinessDate
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

func (dr *DrowsinessRepository) GetDrowsinessStatGroupByHour(from, to *timestamppb.Timestamp, driverUsername *string) ([24]int32, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("drowsiness")
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

	countEachHour := [24]int32{}
	for cur.Next(context.TODO()) {
		var elem models.NumberOfDrowsiness
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		countEachHour[int(elem.ID.Hour)] = elem.Total
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return countEachHour, nil
}
