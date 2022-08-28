# noaa [![GoDoc](https://godoc.org/github.com/icodealot/noaa?status.svg)](https://godoc.org/github.com/icodealot/noaa)

Go package for parts of the weather.gov API. The data provided by weather.gov is in the public domain and covers the continental United States. The service is maintained by the National Weather Service under the umbrella of the National Oceanic and Atmospheric Administration (NOAA). 

Data on various weather.gov API endpoints is measured at different intervals. If a data point is measured hourly then you should take this into account when polling for updates.

## API

This API is currently a minimal subset of what api.weather.gov supports and includes the following:

```go
noaa.Points(lat string, lon string) (points *PointsResponse, err error)
```

```go
noaa.Office(id string) (office *OfficeResponse, err error)
```

```go
noaa.Stations(lat string, lon string) (stations *StationsResponse, err error) {
```

```go
noaa.Forecast(lat string, lon string) (forecast *ForecastResponse, err error) {
```

```go
noaa.GridpointForecast(lat string, lon string) (forecast *GridpointForecastResponse, err error) {
```

```go
noaa.HourlyForecast(lat string, long string) (forecast *HourlyForecastResponse, err error) {
```

For convenience, the ForecastResponse includes a reference to the PointsResponse obtained. In 2017 api.weather.gov was updated with a new REST API that requires multiple calls to obtain the relevant information for the coordinates given by latitude and longitude.

## Setup

Assuming a working `go` toolchain is in place this module can be installed with:

```
go get -u github.com/icodealot/noaa
```

## Examples

There are testable examples in `example_test.go` which can be run using:

```
go test -tags=examples -v
```

A specific example can be run using:

```
go test -tags=examples -run ^ExampleGetChicagoForecast$ -v
```

Here is an example of using the `github.com/icodealot/noaa` module to get forecasted temperatures by day.

```go
package main

import (
	"fmt"
	"github.com/icodealot/noaa"
)

func main() {
	forecast, err := noaa.Forecast("41.837", "-87.685") // Chicago, IL
	if err != nil {
		fmt.Printf("Error getting the forecast: %v", err)
		return
	}
	for _, period := range forecast.Periods {
		fmt.Printf("%-20s ---> %.0f%s\n", period.Name, period.Temperature, period.TemperatureUnit)
	}
}
```

Which will output something like the following:

```bash
This Afternoon       ---> 59F
Tonight              ---> 55F
Tuesday              ---> 67F
Tuesday Night        ---> 47F
Wednesday            ---> 72F
Wednesday Night      ---> 50F
Thursday             ---> 72F
Thursday Night       ---> 51F
Friday               ---> 62F
Friday Night         ---> 36F
Saturday             ---> 52F
Saturday Night       ---> 36F
Sunday               ---> 45F
Sunday Night         ---> 28F
```

Check out the types in `noaa.go` for more details about fields returned by the weather API.
