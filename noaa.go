// Package noaa implements a basic wrapper around api.weather.gov to
// grab HTTP responses to endpoints (i.e.: weather & forecast data)
// by the National Weather Service, an agency of the United States.
package noaa

// Cache used for point lookup to save some HTTP round trips
// key is expected to be PointsResponse.ID
var pointsCache = map[string]*PointsResponse{}

// Points returns a reference to a PointsResponse (cached if appropriate)
// which contains useful noaa endpoints for a given <lat,lon> to use in
// subsequent calls to the api
func Points(lat string, lon string) (points *PointsResponse, err error) {
	endpoint := config.endpointPoints(lat, lon)
	if pointsCache[endpoint] != nil {
		return pointsCache[endpoint], nil
	}
	err = decode(endpoint, &points)
	if err != nil {
		return nil, err
	}
	pointsCache[endpoint] = points
	return
}

// Office returns a reference to a OfficeResponse which contains details
// for a specific forecast office identified by ID
// For example, https://api.weather.gov/offices/LOT (Chicago)
func Office(id string) (office *OfficeResponse, err error) {
	err = decode(config.endpointOffices(id), &office)
	if err != nil {
		return nil, err
	}
	return
}

// Stations returns an array of observation station IDs (urls)
func Stations(lat string, lon string) (stations *StationsResponse, err error) {
	point, err := Points(lat, lon)
	if err != nil {
		return nil, err
	}
	err = decode(point.EndpointObservationStations, &stations)
	if err != nil {
		return nil, err
	}
	return
}

// Forecast returns an array of forecast observations (14 periods and 2/day max)
func Forecast(lat string, lon string) (forecast *ForecastResponse, err error) {
	point, err := Points(lat, lon)
	if err != nil {
		return nil, err
	}
	err = decode(point.EndpointForecast+config.getUnitsQueryParam("?"), &forecast)
	if err != nil {
		return nil, err
	}
	forecast.Point = point
	return
}

// GridpointForecast returns an array of raw forecast data
func GridpointForecast(lat string, long string) (forecast *GridpointForecastResponse, err error) {
	point, err := Points(lat, long)
	if err != nil {
		return nil, err
	}
	err = decode(point.EndpointForecastGridData+config.getUnitsQueryParam("?"), &forecast)
	if err != nil {
		return nil, err
	}
	forecast.Point = point
	return forecast, nil
}

// HourlyForecast returns an array of raw hourly forecast data
func HourlyForecast(lat string, long string) (forecast *HourlyForecastResponse, err error) {
	point, err := Points(lat, long)
	if err != nil {
		return nil, err
	}
	err = decode(point.EndpointForecastHourly+config.getUnitsQueryParam("?"), &forecast)
	if err != nil {
		return nil, err
	}
	forecast.Point = point
	return forecast, nil
}
