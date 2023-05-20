// Package noaa implements a basic wrapper around api.weather.gov to
// grab HTTP responses to endpoints (i.e.: weather & forecast data)
// by the National Weather Service, an agency of the United States.
package noaa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

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
	req.Header.Add("Accept", config.Accept)
	req.Header.Add("User-Agent", config.UserAgent)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%d %s", res.StatusCode, res.Status))
	}

	return res, nil
}

// Points returns a set of useful endpoints for a given <lat,lon>
// or returns a cached object if appropriate
func Points(lat string, lon string) (points *PointsResponse, err error) {
	endpoint := fmt.Sprintf("%s/points/%s,%s", config.BaseURL, lat, lon)
	if pointsCache[endpoint] != nil {
		return pointsCache[endpoint], nil
	}
	res, err := apiCall(endpoint)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&points); err != nil {
		return nil, err
	}
	pointsCache[endpoint] = points
	return points, nil
}

// Office returns details for a specific office identified by its ID
// For example, https://api.weather.gov/offices/LOT (Chicago)
func Office(id string) (office *OfficeResponse, err error) {
	endpoint := fmt.Sprintf("%s/offices/%s", config.BaseURL, id)

	res, err := apiCall(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&office); err != nil {
		return nil, err
	}
	return office, nil
}

// Stations returns an array of observation station IDs (urls)
func Stations(lat string, lon string) (stations *StationsResponse, err error) {
	point, err := Points(lat, lon)
	if err != nil {
		return nil, err
	}
	res, err := apiCall(point.EndpointObservationStations)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&stations); err != nil {
		return nil, err
	}
	return stations, nil
}

// Forecast returns an array of forecast observations (14 periods and 2/day max)
func Forecast(lat string, lon string) (forecast *ForecastResponse, err error) {
	query := ""
	point, err := Points(lat, lon)
	if err != nil {
		return nil, err
	}
	if config.Units != "" {
		query = "?units=" + config.Units
	}
	res, err := apiCall(point.EndpointForecast + query)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&forecast); err != nil {
		return nil, err
	}
	forecast.Point = point
	return forecast, nil
}

// GridpointForecast returns an array of raw forecast data
func GridpointForecast(lat string, long string) (forecast *GridpointForecastResponse, err error) {
	query := ""
	point, err := Points(lat, long)
	if err != nil {
		return nil, err
	}
	if config.Units != "" {
		query = "?units=" + config.Units
	}
	res, err := apiCall(point.EndpointForecastGridData + query)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&forecast); err != nil {
		return nil, err
	}
	forecast.Point = point
	return forecast, nil
}

// HourlyForecast returns an array of raw hourly forecast data
func HourlyForecast(lat string, long string) (forecast *HourlyForecastResponse, err error) {
	query := ""
	point, err := Points(lat, long)
	if err != nil {
		return nil, err
	}
	if config.Units != "" {
		query = "?units=" + config.Units
	}
	res, err := apiCall(point.EndpointForecastHourly + query)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&forecast); err != nil {
		return nil, err
	}
	forecast.Point = point
	return forecast, nil
}
