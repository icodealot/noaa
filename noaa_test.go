package noaa_test

import (
	"testing"

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

func TestChicagoOffice(t *testing.T) {
	office, err := noaa.Office("LOT")
	if office != nil && err == nil {
		if office.Name == "Chicago, IL" {
			return
		}
	}
	t.Error("noaa.Office(\"LOT\") should return valid office information.")
}
