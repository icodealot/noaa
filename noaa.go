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
	EndpointForecastGridData    string `json:"forecastGridData"`
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
		ID              int32   `json:"number"`
		Name            string  `json:"name"`
		StartTime       string  `json:"startTime"`
		EndTime         string  `json:"endTime"`
		IsDaytime       bool    `json:"isDaytime"`
		Temperature     float64 `json:"temperature"`
		TemperatureUnit string  `json:"temperatureUnit"`
		WindSpeed       string  `json:"windSpeed"`
		WindDirection   string  `json:"windDirection"`
		Summary         string  `json:"shortForecast"`
		Details         string  `json:"detailedForecast"`
	} `json:"periods"`
	Point *PointsResponse
}

// GridpointForecastResponse holds the JSON values from /gridpoints/<cwa>/<x,y>"
// See https://weather-gov.github.io/api/gridpoints for information.
type GridpointForecastResponse struct {
	// capture data from the forecast
	Updated   string `json:"updateTime"`
	Elevation struct {
		Value float64 `json:"value"`
		Units string  `json:"unitCode"`
	} `json:"elevation"`
	Weather struct {
		Values []struct {
			ValidTime string `json:"validTime"` // ISO 8601 time interval, e.g. 2019-07-04T18:00:00+00:00/PT3H
			Value     []struct {
				Coverage  string `json:"coverage"`
				Weather   string `json:"weather"`
				Intensity string `json:"intensity"`
			} `json:"value"`
		} `json:"values"`
	} `json:"weather"`
	Hazards struct {
		Values []struct {
			ValidTime string `json:"validTime"` // ISO 8601 time interval, e.g. 2019-07-04T18:00:00+00:00/PT3H
			Value     []struct {
				Phenomenon   string `json:"phenomenon"`
				Significance string `json:"significance"`
				EventNumber  int32  `json:"event_number"`
			} `json:"value"`
		} `json:"values"`
	} `json:"hazards"`
	Temperature                      GridpointForecastTimeSeries `json:"temperature"`
	Dewpoint                         GridpointForecastTimeSeries `json:"dewpoint"`
	MaxTemperature                   GridpointForecastTimeSeries `json:"maxTemperature"`
	MinTemperature                   GridpointForecastTimeSeries `json:"minTemperature"`
	RelativeHumidity                 GridpointForecastTimeSeries `json:"relativeHumidity"`
	ApparentTemperature              GridpointForecastTimeSeries `json:"apparentTemperature"`
	HeatIndex                        GridpointForecastTimeSeries `json:"heatIndex"`
	WindChill                        GridpointForecastTimeSeries `json:"windChill"`
	SkyCover                         GridpointForecastTimeSeries `json:"skyCover"`
	WindDirection                    GridpointForecastTimeSeries `json:"windDirection"`
	WindSpeed                        GridpointForecastTimeSeries `json:"windSpeed"`
	WindGust                         GridpointForecastTimeSeries `json:"windGust"`
	ProbabilityOfPrecipitation       GridpointForecastTimeSeries `json:"probabilityOfPrecipitation"`
	QuantitativePrecipitation        GridpointForecastTimeSeries `json:"quantitativePrecipitation"`
	iceAccumulation                  GridpointForecastTimeSeries `json:"iceAccumulation"`
	snowfallAmount                   GridpointForecastTimeSeries `json:"snowfallAmount"`
	snowLevel                        GridpointForecastTimeSeries `json:"snowLevel"`
	ceilingHeight                    GridpointForecastTimeSeries `json:"ceilingHeight"`
	visibility                       GridpointForecastTimeSeries `json:"visibility"`
	transportWindSpeed               GridpointForecastTimeSeries `json:"transportWindSpeed"`
	transportWindDirection           GridpointForecastTimeSeries `json:"transportWindDirection"`
	mixingHeight                     GridpointForecastTimeSeries `json:"mixingHeight"`
	hainesIndex                      GridpointForecastTimeSeries `json:"hainesIndex"`
	lightningActivityLevel           GridpointForecastTimeSeries `json:"lightningActivityLevel"`
	twentyFootWindSpeed              GridpointForecastTimeSeries `json:"twentyFootWindSpeed"`
	twentyFootWindDirection          GridpointForecastTimeSeries `json:"twentyFootWindDirection"`
	waveHeight                       GridpointForecastTimeSeries `json:"waveHeight"`
	wavePeriod                       GridpointForecastTimeSeries `json:"wavePeriod"`
	waveDirection                    GridpointForecastTimeSeries `json:"waveDirection"`
	primarySwellHeight               GridpointForecastTimeSeries `json:"primarySwellHeight"`
	primarySwellDirection            GridpointForecastTimeSeries `json:"primarySwellDirection"`
	secondarySwellHeight             GridpointForecastTimeSeries `json:"secondarySwellHeight"`
	secondarySwellDirection          GridpointForecastTimeSeries `json:"secondarySwellDirection"`
	wavePeriod2                      GridpointForecastTimeSeries `json:"wavePeriod2"`
	windWaveHeight                   GridpointForecastTimeSeries `json:"windWaveHeight"`
	dispersionIndex                  GridpointForecastTimeSeries `json:"dispersionIndex"`
	pressure                         GridpointForecastTimeSeries `json:"pressure"`
	probabilityOfTropicalStormWinds  GridpointForecastTimeSeries `json:"probabilityOfTropicalStormWinds"`
	probabilityOfHurricaneWinds      GridpointForecastTimeSeries `json:"probabilityOfHurricaneWinds"`
	potentialOf15mphWinds            GridpointForecastTimeSeries `json:"potentialOf15mphWinds"`
	potentialOf25mphWinds            GridpointForecastTimeSeries `json:"potentialOf25mphWinds"`
	potentialOf35mphWinds            GridpointForecastTimeSeries `json:"potentialOf35mphWinds"`
	potentialOf45mphWinds            GridpointForecastTimeSeries `json:"potentialOf45mphWinds"`
	potentialOf20mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf20mphWindGusts"`
	potentialOf30mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf30mphWindGusts"`
	potentialOf40mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf40mphWindGusts"`
	potentialOf50mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf50mphWindGusts"`
	potentialOf60mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf60mphWindGusts"`
	grasslandFireDangerIndex         GridpointForecastTimeSeries `json:"grasslandFireDangerIndex"`
	probabilityOfThunder             GridpointForecastTimeSeries `json:"probabilityOfThunder"`
	davisStabilityIndex              GridpointForecastTimeSeries `json:"davisStabilityIndex"`
	atmosphericDispersionIndex       GridpointForecastTimeSeries `json:"atmosphericDispersionIndex"`
	lowVisibilityOccurrenceRiskIndex GridpointForecastTimeSeries `json:"lowVisibilityOccurrenceRiskIndex"`
	stability                        GridpointForecastTimeSeries `json:"stability"`
	redFlagThreatIndex               GridpointForecastTimeSeries `json:"redFlagThreatIndex"`
	Point                            *PointsResponse
}

// GridpointForecastTimeSeries holds a series of data from a gridpoint forecast
type GridpointForecastTimeSeries struct {
	Uom    string `json:"uom"` // Unit of Measure
	Values []struct {
		ValidTime string  `json:"validTime"` // ISO 8601 time interval, e.g. 2019-07-04T18:00:00+00:00/PT3H
		Value     float64 `json:value"`
	} `json:"values"`
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

	if res.StatusCode == 404 {
		defer res.Body.Close()
		return nil, errors.New("404: data not found for -> " + endpoint)
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

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&points); err != nil {
		return nil, err
	}
	pointsCache[endpoint] = points
	return points, nil
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
	point, err := Points(lat, lon)
	if err != nil {
		return nil, err
	}
	res, err := apiCall(point.EndpointForecast)
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
	point, err := Points(lat, long)
	if err != nil {
		return nil, err
	}
	res, err := apiCall(point.EndpointForecastGridData)
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
