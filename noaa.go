// Package noaa implements a basic wrapper around api.weather.gov to
// grab HTTP responses to endpoints (i.e.: weather & forecast data)
// by the National Weather Service, an agency of the United States.
package noaa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Constant values for the weather.gov REST API
const (
	API       = "https://api.weather.gov"
	APIKey    = "github.com/icodealot/noaa" // See auth docs at weather.gov
	APIAccept = "application/ld+json"       // Changes may affect struct mappings below
)

// PointsResponse holds the JSON values from /points/<lat,lon>
type PointsResponse struct {
	ID                          string `json:"@id"`
	CWA                         string `json:"cwa"`
	Office                      string `json:"forecastOffice"`
	GridX                       int64  `json:"gridX"`
	GridY                       int64  `json:"gridY"`
	EndpointForecast            string `json:"forecast"`
	EndpointForecastHourly      string `json:"forecastHourly"`
	EndpointObservationStations string `json:"observationStations"`
	Timezone                    string `json:"timeZone"`
	RadarStation                string `json:"radarStation"`
}

// StationsResponse holds the JSON values from /points/<lat,lon>/stations
type StationsResponse struct {
	Stations []string `json:"observationStations"`
}

// ForecastResponse holds the JSON values from /gridpoints/<cwa>/<x,y>/forecast"
type ForecastResponse struct {
	// capture data from the forecast
	Updated   string `json:"updated"`
	Units     string `json:"units"`
	Elevation struct {
		Value float64 `json:"value"`
		Units string  `json:"unitCode"`
	} `json:"elevation"`
	Periods []struct {
		ID             int32   `json:"number"`
		Name           string  `json:"name"`
		StartTime      string  `json:"startTime"`
		EndTime        string  `json:"endTime"`
		IsDaytime      bool    `json:"isDaytime"`
		Temperature    float64 `json:"temperature"`
		TemeratureUnit string  `json:"temperatureUnit"`
		WindSpeed      string  `json:"windSpeed"`
		WindDirection  string  `json:"windDirection"`
		Summary        string  `json:"shortForecast"`
		Details        string  `json:"detailedForecast"`
	} `json:"periods"`
	Point *PointsResponse
}

// Cache used for point lookup to save some HTTP round trips
// key is expected to be PointsResponse.ID
var pointsCache = map[string]*PointsResponse{}

// Call the weather.gov API. We could just use http.Get() but
// since we need to include some custom header values this helps.
func apiCall(endpoint string) (res *http.Response, err error) {
	endpoint = strings.Replace(endpoint, "http://", "https://", -1)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", APIAccept)
	req.Header.Add("User-Agent", APIKey) // See http://www.weather.gov/documentation/services-web-api

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Points returns a set of useful endpoints for a given <lat,lon>
// or returns a cached object if appropriate
func Points(lat string, lon string) (points *PointsResponse, err error) {
	endpoint := fmt.Sprintf("%s/points/%s,%s", API, lat, lon)
	if pointsCache[endpoint] != nil {
		return pointsCache[endpoint], nil
	}
	res, err := apiCall(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var jsonResponse PointsResponse
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&jsonResponse); err != nil {
		return nil, err
	}
	pointsCache[endpoint] = &jsonResponse
	return &jsonResponse, nil
}

// Stations returns an array of observation station IDs (urls)
func Stations(lat string, lon string) ([]string, error) {
	point, err := Points(lat, lon)
	if err != nil {
		return nil, err
	}
	res, err := apiCall(point.EndpointObservationStations)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var jsonResponse *StationsResponse

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&jsonResponse); err != nil {
		return nil, err
	}
	return jsonResponse.Stations, nil
}

// Forecast returns an array of forecast observations (14 periods and 2/day max)
func Forecast(lat string, lon string) (forecast *ForecastResponse, err error) {
	point, err := Points(lat, lon)
	if err != nil {
		return nil, err
	}
	res, err := apiCall(point.EndpointForecast)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var jsonResponse *ForecastResponse

	// b, err := ioutil.ReadAll(res.Body)
	// fmt.Printf("%s\n", b)

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&jsonResponse); err != nil {
		return nil, err
	}
	jsonResponse.Point = point
	return jsonResponse, nil
}
