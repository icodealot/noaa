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
	if err != nil {
		t.Error("noaa.Forecast() should return valid data for Chicago.")
	}
	if forecast.Units == "si" {
		return
	}
	t.Error("noaa.Forecast() should return valid data for Chicago in Metric.")
}

func TestUSUnits(t *testing.T) {
	noaa.SetUnits("us")
	forecast, err := noaa.Forecast("41.837", "-87.685")
	if err != nil {
		t.Error("noaa.Forecast() should return valid data for Chicago.")
	}
	if forecast.Units == "us" {
		return
	}
	t.Error("noaa.Forecast() should return valid data for Chicago in standard units.")
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
	client := &http.Client{
		Timeout: time.Millisecond * 1,
	}
	orig := noaa.SetClient(client)
	if orig != http.DefaultClient {
		t.Error("expected http.DefaultClient but got something else.")
	}

	_, err := noaa.HourlyForecast("41.837", "-87.685")
	if err == nil {
		t.Error("expected request to time out (1 millisecond is too brief for any request)")
	}

	myClient := noaa.SetClient(orig)
	if myClient != client {
		t.Error("expected client, but got something else.")
	}

	_, err = noaa.HourlyForecast("41.837", "-87.685")
	if err != nil {
		t.Error("noaa.HourlyForecast() should return valid data.")
	}
}
