// Package noaa implements a basic wrapper around api.weather.gov to
// grab HTTP responses to endpoints (i.e.: weather & forecast data)
// by the National Weather Service, an agency of the United States.
package noaa

import "fmt"

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
	updateForecastPeriods(forecast.Periods)
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
	updateForecastPeriods(forecast.Periods)
	return forecast, nil
}

// Using the quantitative value feature flags to enable QV responses
// causes the noaa api to ignore the requested unit types. This also
// populates fields that were previously populated for backward
// compatibility. This is necessary because quantitative values replace
// deprecated fields with a nested object. See: QuantitativeValue.
// These are nice to have but may be deprecated in the future.
func updateForecastPeriods(periods []ForecastResponsePeriod) {
	for i, period := range periods {
		updateTemperature(&period)
		updateWindSpeed(&period)
		periods[i] = period
	}
}

// See: updateForecastPeriods
func updateTemperature(period *ForecastResponsePeriod) {
	wmoUnitCode := period.QuantitativeTemperature.UnitCode
	period.Temperature = period.QuantitativeTemperature.Value
	if config.Units == "si" {
		period.TemperatureUnit = "C"
		if wmoUnitCode != "wmoUnit:degC" {
			// assume its degrees F so convert it accordingly
			period.Temperature = (5.0 / 9.0) * (period.Temperature - 32)
		}
	} else {
		period.TemperatureUnit = "F"
		if wmoUnitCode == "wmoUnit:degC" {
			period.Temperature = ((9.0 / 5.0) * period.Temperature) + 32
		}
	}
}

const (
	KilometersPerMile = 1.60934
	MilesPerKilometer = 0.62137
)

// See: updateForecastPeriods
func updateWindSpeed(period *ForecastResponsePeriod) {
	wmoUnitCode := period.QuantitativeWindSpeed.UnitCode
	min := period.QuantitativeWindSpeed.MinValue
	max := period.QuantitativeWindSpeed.MaxValue
	value := period.QuantitativeWindSpeed.Value
	units := ""

	if config.Units == "si" {
		units = "km/h"
		if wmoUnitCode != "wmoUnit:km_h-1" {
			// assume its mph so convert it accordingly
			min *= KilometersPerMile
			max *= KilometersPerMile
			value *= KilometersPerMile
		}
	} else {
		units = "mph"
		if wmoUnitCode == "wmoUnit:km_h-1" {
			// assume its kmh so convert it accordingly
			min *= MilesPerKilometer
			max *= MilesPerKilometer
			value *= MilesPerKilometer
		}
	}

	// replicates legacy api behavior but using quantitative values
	if min == 0.0 && max == 0.0 {
		period.WindSpeed = fmt.Sprintf("%.0f %s", value, units)
	} else {
		period.WindSpeed = fmt.Sprintf("%.0f to %.0f %s", min, max, units)
	}
}
