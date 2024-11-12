//go:build examples
// +build examples

// Examples can be run with `go test -tags=examples -v` and do not necessarily
// require access to the weather.gov API. Examples demonstrate various uses of
// the client and how to set config values.
package noaa_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/icodealot/noaa"
)

func ExampleSetConfig() {

	// Cleanup global state before each example
	beforeEachExample()

	// When replacing the entire Config, three fields are required at a minimum.
	// In general SetUnits() and SetUserAgent() are probably more useful.
	noaa.SetConfig(noaa.Config{
		BaseURL:   "https://unplanned-hostname-change.com",
		UserAgent: "(github.com/icodealot/noaa test user-agent)",
		Accept:    "application/ld+json",
	})

	// The current configuration can be retrieved as follows:
	config := noaa.GetConfig()

	fmt.Println("Config changed to:", config.BaseURL, config.UserAgent)

	// Output:
	// Config changed to: https://unplanned-hostname-change.com (github.com/icodealot/noaa test user-agent)
}

func ExampleSetUserAgent() {

	// Cleanup global state before each example
	beforeEachExample()

	// Set the user-agent field for your own application. For more information
	// See Authentication located at: https://www.weather.gov/documentation/services-web-api
	noaa.SetUserAgent("(myweatherapp.com, contact@myweatherapp.com)")

	// Get the current configuration:
	config := noaa.GetConfig()

	fmt.Println("User-Agent should now be:", config.UserAgent)

	// Output:
	// User-Agent should now be: (myweatherapp.com, contact@myweatherapp.com)
}

func ExampleSetUnits() {

	// Cleanup global state before each example
	beforeEachExample()

	// Set the units.
	// Units can be set to "us" or "si" and otherwise, blank "" defaults to US units.
	noaa.SetUnits("si")

	// Get the current configuration:
	config := noaa.GetConfig()

	fmt.Println("Units should now be:", config.Units)

	// Output:
	// Units should now be: si
}

func ExampleGetChicagoForecast() {

	// Cleanup global state before each example
	beforeEachExample()

	// Get the forecast for Chicago by lat/lon
	forecast, err := noaa.Forecast("41.837", "-87.685")
	if err != nil {
		fmt.Printf("Error getting the forecast: %v", err)
		return
	}
	for _, period := range forecast.Periods {
		log.Printf("%-20s ---> Windspeed: %-15s Temperature: %.0f%s\n", period.Name, period.WindSpeed, period.Temperature, period.TemperatureUnit)
	}

	fmt.Println("Success!")

	// Output:
	// Success!
}

func ExampleGetChicagoForecastWithMetricUnits() {

	// Cleanup global state before each example
	beforeEachExample()

	noaa.SetUnits("si")

	// Get the forecast for Chicago by lat/lon
	forecast, err := noaa.Forecast("41.837", "-87.685")
	if err != nil || forecast == nil {
		fmt.Printf("Error getting the forecast: %v", err)
		return
	}
	for _, period := range forecast.Periods {
		log.Printf("%-20s ---> Windspeed: %-15s Temperature: %.0f%s\n", period.Name, period.WindSpeed, period.Temperature, period.TemperatureUnit)
	}

	fmt.Println("Success!")

	// Output:
	// Success!
}

func ExampleGetChicagoHourlyForecast() {

	// Cleanup global state before each example
	beforeEachExample()

	// Get the hourly forecast for Chicago by lat/lon
	response, err := noaa.HourlyForecast("41.837", "-87.685")
	if err != nil {
		fmt.Printf("Error getting the forecast: %v", err)
		return
	}
	for _, period := range response.Periods {
		log.Printf("at %s ... it will be %.0f%s\n", period.StartTime, period.Temperature, period.TemperatureUnit)
	}

	fmt.Println("Success!")

	// Output:
	// Success!
}

func ExampleGetChicagoGridpointForecast() {

	// Cleanup global state before each example
	beforeEachExample()

	// Get the gridpoint forecast for Chicago by lat/lon
	response, err := noaa.GridpointForecast("41.837", "-87.685")
	if err != nil {
		fmt.Printf("Error getting the gridpoint forecast: %v", err)
		return
	}
	log.Printf("Gridpoint forecast:\n%+v\n", response)

	fmt.Println("Success!")

	// Output:
	// Success!
}

func ExampleGetChicagoWeatherStations() {

	// Cleanup global state before each example
	beforeEachExample()

	// Get the hourly forecast for Chicago by lat/lon
	response, err := noaa.Stations("41.837", "-87.685")
	if err != nil {
		fmt.Printf("Error getting the stations: %v", err)
		return
	}
	for _, station := range response.Stations {
		log.Printf("Weather station: %s\n", station)
	}

	fmt.Println("Success!")

	// Output:
	// Success!
}

func ExampleSetClient() {

	// Cleanup global state before each example
	beforeEachExample()

	// Use a client that can't successfully make a request.
	noaa.SetClient(&http.Client{
		Timeout: time.Millisecond * 1,
	})

	// Get the hourly forecast for Chicago by lat/lon
	_, err := noaa.Stations("41.837", "-87.685")
	if err != nil {
		fmt.Printf("Error getting the stations: %v", err)
		return
	}

	fmt.Println("Success!")

	// Output:
	// Error getting the stations: Get "https://api.weather.gov/gridpoints/LOT/74,71/stations": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
}

// beforeEachExample is used to clean up the global state of the noaa client
// which is necessary because some global state is set at the module level
func beforeEachExample() {
	// Reset the config for this example (cleanup from above examples)
	noaa.SetConfig(noaa.GetDefaultConfig())
}
