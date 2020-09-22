# noaa [![GoDoc](https://godoc.org/github.com/icodealot/noaa?status.svg)](https://godoc.org/github.com/icodealot/noaa)

Go package for parts of the weather.gov API. The data provided by weather.gov is in the public domain and covers the continental United States. The service is maintained by the National Weather Service under the umbrella of the National Oceanic and Atmospheric Administration (NOAA). 

Data on various weather.gov API endpoints is measured at different intervals. If a data point is measured hourly then you should take this into account when polling for updates.

## API

This API is currenly a minimal subset of what api.weather.gov supports and includes the following:

```go
noaa.Points(lat string, lon string) (points *PointsResponse, err error)
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

For convenience, the ForecastResponse includes a reference to the PointsResponse obtained. In 2017 api.weather.gov was updated with a new REST API that requires multiple calls to obtain the relevant information for the coordinates given by latitude and longitude.

## Example

Here is an example of using the `github.com/icodealot/noaa` package for Go to pull the forecasted temperatures by day.

```go
package main

import (
	"fmt"
	"github.com/icodealot/noaa"
)

func main() {
	forecast, err := noaa.Forecast("30.5835", "-97.8575")
	if err != nil {
		fmt.Printf("Error getting the forecast: %v", err)
		return
	}
	for _, period := range forecast.Periods {
		fmt.Printf("%-20s ---> %.0f\n", period.Name, period.Temperature)
	}
}
```

Which will output something like the following:

```bash
This Afternoon       ---> 59
Tonight              ---> 55
Tuesday              ---> 67
Tuesday Night        ---> 47
Wednesday            ---> 72
Wednesday Night      ---> 50
Thursday             ---> 72
Thursday Night       ---> 51
Friday               ---> 62
Friday Night         ---> 36
Saturday             ---> 52
Saturday Night       ---> 36
Sunday               ---> 45
Sunday Night         ---> 28
```

Check out the code for more details. This package is still super unfinished and pretty unstable with respect to changes so if you want to insulate yourself then please fork the project as needed.
