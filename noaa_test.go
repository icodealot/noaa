//go:build !examples
// +build !examples

// Unit tests can be run with `go test -v` and require access to the API. Many of
// these tests are actually integration tests that call the weather.gov API and
// parse responses accordingly to confirm expected responses are returned.
//
// Thus, in the future if weather.gov changes the endpoints or responses, these
// tests should alert users of this wrapper SDK accordingly.
package noaa_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/icodealot/noaa"
)

func TestBlank(t *testing.T) {
	point, err := noaa.Points("", "")
	if point == nil && err != nil {
		return
	}
	t.Error("noaa.Points() should return a 404 error for a blank lat, lon.")
}

func TestBlankLat(t *testing.T) {
	point, err := noaa.Points("", "-147.7390417")
	if point == nil && err != nil {
		return
	}
	t.Error("noaa.Points() should return a 404 error for a blank lat.")
}

func TestBlankLon(t *testing.T) {
	point, err := noaa.Points("64.828421", "")
	if point == nil && err != nil {
		return
	}
	t.Error("noaa.Points() should return a 404 error for a blank lon.")
}

func TestZero(t *testing.T) {
	point, err := noaa.Points("0", "0")
	if point == nil && err != nil {
		return
	}
	t.Error("noaa.Points() should return a 404 error for a zero lat, lon.")
}

func TestInternational(t *testing.T) {
	point, err := noaa.Points("48.85660", "2.3522") // Paris, France
	if point == nil && err != nil {
		return
	}
	t.Error("noaa.Points() should return a 404 error for lat, lon outside the U.S. territories.")
}

func TestAlaska(t *testing.T) {
	point, err := noaa.Points("64.828421", "-147.7390417")
	if point != nil && err == nil {
		return
	}
	t.Error("noaa.Points() should return valid points for parts of Alaska.")
}

func TestMetricUnits(t *testing.T) {
	noaa.SetUnits("si")
	forecast, err := noaa.Forecast("41.837", "-87.685")
	if err != nil || forecast == nil {
		t.Error("noaa.Forecast() should return valid data for Chicago.")
		return
	}
	if forecast.Units != "si" {
		t.Error("noaa.Forecast() should return valid data for Chicago in metric.")
	}
}

func TestUSUnits(t *testing.T) {
	noaa.SetUnits("us")
	forecast, err := noaa.Forecast("41.837", "-87.685")
	if err != nil {
		t.Error("noaa.Forecast() should return valid data for Chicago.")
	}
	if forecast.Units != "us" {
		t.Error("noaa.Forecast() should return valid data for Chicago in standard units.")
	}
}

func TestChicagoOffice(t *testing.T) {
	office, err := noaa.Office("LOT")
	if office != nil && err == nil {
		if office.Name == "Chicago, IL" {
			return
		}
	}
	t.Error("noaa.Office(\"LOT\") should return valid office information.")
}

func TestChicagoHourly(t *testing.T) {
	hourly, err := noaa.HourlyForecast("41.837", "-87.685")
	if err != nil {
		t.Error("noaa.HourlyForecast() should return valid data for Chicago.")
	}
	if len(hourly.Periods) == 0 {
		t.Error("expected at least one period")
	}
}

func TestSetClient(t *testing.T) {
	// Intentionally create a client with an absurd timeout value.
	client := &http.Client{
		Timeout: time.Millisecond,
	}

	// Don't set the client, to ensure this still works as normal.
	_, err := noaa.HourlyForecast("41.837", "-87.685")
	if err != nil {
		t.Errorf("should have successfully returned a result instead of this error: %s", err)
	}

	// Set the client to test this feature.
	noaa.SetClient(client)

	// See if we can make a (failing) request with this.
	_, err = noaa.HourlyForecast("41.837", "-87.685")
	if err == nil {
		t.Error("should have failed the request, 1 millisecond is too short a timeout to make the request")
	}

	// Test that setting to nil returns to http.DefaultClient.
	noaa.SetClient(nil)
	_, err = noaa.HourlyForecast("41.837", "-87.685")
	if err != nil {
		t.Errorf("should have successfully returned a result instead of this error: %s", err)
	}
}
