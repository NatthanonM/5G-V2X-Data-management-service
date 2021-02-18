package services

import (
	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GoogleService struct {
	config *config.Config
}

func NewGoogleService(config *config.Config) *GoogleService {
	return &GoogleService{
		config: config,
	}
}

// ReverseGeocoding ...
func (gs *GoogleService) ReverseGeocoding(lat, lng float64) (*string, error) {
	url := fmt.Sprintf(`https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&result_type=route&key=%s`, lat, lng, gs.config.GoogleAPIKey)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:")
	fmt.Println(string(body))

	var reverseGeocode models.ReverseGeocode
	if err := json.Unmarshal(body, &reverseGeocode); err != nil {
		log.Println(err)
		return nil, err
	}
	roadName := "Unnamed Road"
	if len(reverseGeocode.Results) != 0 {
		roadName = reverseGeocode.Results[0].AddressComponents[0].LongName
	}

	return &roadName, nil
}
