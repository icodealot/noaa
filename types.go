package noaa

import (
	"time"
)

// QuantitativeValue is available for various statistics and can be
// enabled with an optional request header to the noaa API. In the
// future it is expected at that QV will replace single values such
// as measurements for Temperature
type QuantitativeValue struct {
	Value          float64 `json:"value"`
	MaxValue       float64 `json:"maxValue"`
	MinValue       float64 `json:"minValue"`
	UnitCode       string  `json:"unitCode"`
	QualityControl string  `json:"qualityControl"`
}

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

// OfficeAddress holds the JSON values for the address of an OfficeResponse
type OfficeAddress struct {
	Type          string `json:"@type"`
	StreetAddress string `json:"streetAddress"`
	Locality      string `json:"addressLocality"`
	Region        string `json:"addressRegion"`
	PostalCode    string `json:"postalCode"`
}

// OfficeResponse holds the JSON values from /offices/<id>
type OfficeResponse struct {
	Type                        string        `json:"@type"`
	URI                         string        `json:"@id"`
	ID                          string        `json:"id"`
	Name                        string        `json:"name"`
	Address                     OfficeAddress `json:"address"`
	Telephone                   string        `json:"telephone"`
	FaxNumber                   string        `json:"faxNumber"`
	Email                       string        `json:"email"`
	SameAs                      string        `json:"sameAs"`
	NWSRegion                   string        `json:"nwsRegion"`
	ParentOrganization          string        `json:"parentOrganization"`
	ResponsibleCounties         []string      `json:"responsibleCounties"`
	ResponsibleForecastZones    []string      `json:"responsibleForecastZones"`
	ResponsibleFireZones        []string      `json:"responsibleFireZones"`
	ApprovedObservationStations []string      `json:"approvedObservationStations"`
}

// StationsResponse holds the JSON values from /points/<lat,lon>/stations
type StationsResponse struct {
	Stations []string `json:"observationStations"`
}

// ForecastElevation holds the JSON values for a forecast response's elevation.
type ForecastElevation struct {
	Value float64 `json:"value"`
	Units string  `json:"unitCode"`
}

// ForecastResponsePeriod holds the JSON values for a period within a forecast response.
type ForecastResponsePeriod struct {
	ID               int32   `json:"number"`
	Name             string  `json:"name"`
	StartTime        string  `json:"startTime"`
	EndTime          string  `json:"endTime"`
	IsDaytime        bool    `json:"isDaytime"`
	Temperature      float64 // preserved for legacy compatibility, may be deprecated in the future
	TemperatureUnit  string  // preserved for legacy compatibility, may be deprecated in the future
	TemperatureTrend string  `json:"temperatureTrend"`
	WindSpeed        string  // preserved for legacy compatibility, may be deprecated in the future
	WindDirection    string  `json:"windDirection"`
	Icon             string  `json:"icon"`
	Summary          string  `json:"shortForecast"`
	Details          string  `json:"detailedForecast"`

	QuantitativeProbability      QuantitativeValue `json:"probabilityOfPrecipitation"`
	QuantitativeDewpoint         QuantitativeValue `json:"dewpoint"`
	QuantitativeRelativeHumidity QuantitativeValue `json:"relativeHumidity"`
	QuantitativeTemperature      QuantitativeValue `json:"temperature"`
	QuantitativeWindSpeed        QuantitativeValue `json:"windSpeed"`
	QuantitativeWindGust         QuantitativeValue `json:"windGust"`
}

// ForecastResponsePeriodHourly provides the JSON value for a period within an hourly forecast.
type ForecastResponsePeriodHourly = ForecastResponsePeriod

// ForecastResponse holds the JSON values from /gridpoints/<cwa>/<x,y>/forecast"
type ForecastResponse struct {
	Updated   string                   `json:"updated"`
	Units     string                   `json:"units"`
	Elevation ForecastElevation        `json:"elevation"`
	Periods   []ForecastResponsePeriod `json:"periods"`
	Point     *PointsResponse
}

// WeatherValueItem holds the JSON values for a weather.values[x].value.
type WeatherValueItem struct {
	Coverage  string `json:"coverage"`
	Weather   string `json:"weather"`
	Intensity string `json:"intensity"`
}

// WeatherValue holds the JSON value for a weather.values[x] value.
type WeatherValue struct {
	ValidTime string             `json:"validTime"` // ISO 8601 time interval, e.g. 2019-07-04T18:00:00+00:00/PT3H
	Value     []WeatherValueItem `json:"value"`
}

// Weather holds the JSON value for the weather object.
type Weather struct {
	Values []WeatherValue `json:"values"`
}

// HazardValueItem holds a value item from a GridpointForecastResponse's
// hazard.values[x].value[x].
type HazardValueItem struct {
	Phenomenon   string `json:"phenomenon"`
	Significance string `json:"significance"`
	EventNumber  int32  `json:"event_number"`
}

// HazardValue holds a hazard value from a GridpointForecastResponse's
// hazard.values[x].
type HazardValue struct {
	ValidTime string            `json:"validTime"` // ISO 8601 time interval, e.g. 2019-07-04T18:00:00+00:00/PT3H
	Value     []HazardValueItem `json:"value"`
}

// Hazard holds a slice of HazardValue items from a GridpointForecastResponse hazards
type Hazard struct {
	Values []HazardValue `json:"values"`
}

// HourlyForecastResponse holds the JSON values for the hourly forecast.
type HourlyForecastResponse struct {
	Updated           string                         `json:"updated"`
	Units             string                         `json:"units"`
	ForecastGenerator string                         `json:"forecastGenerator"`
	GeneratedAt       string                         `json:"generatedAt"`
	UpdateTime        string                         `json:"updateTime"`
	ValidTimes        string                         `json:"validTimes"`
	Periods           []ForecastResponsePeriodHourly `json:"periods"`
	Point             *PointsResponse
}

// GridpointForecastResponse holds the JSON values from /gridpoints/<cwa>/<x,y>"
// See https://weather-gov.github.io/api/gridpoints for information.
type GridpointForecastResponse struct {
	Updated                          string                      `json:"updateTime"`
	Elevation                        ForecastElevation           `json:"elevation"`
	Weather                          Weather                     `json:"weather"`
	Hazards                          Hazard                      `json:"hazards"`
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
	IceAccumulation                  GridpointForecastTimeSeries `json:"iceAccumulation"`
	SnowfallAmount                   GridpointForecastTimeSeries `json:"snowfallAmount"`
	SnowLevel                        GridpointForecastTimeSeries `json:"snowLevel"`
	CeilingHeight                    GridpointForecastTimeSeries `json:"ceilingHeight"`
	Visibility                       GridpointForecastTimeSeries `json:"visibility"`
	TransportWindSpeed               GridpointForecastTimeSeries `json:"transportWindSpeed"`
	TransportWindDirection           GridpointForecastTimeSeries `json:"transportWindDirection"`
	MixingHeight                     GridpointForecastTimeSeries `json:"mixingHeight"`
	HainesIndex                      GridpointForecastTimeSeries `json:"hainesIndex"`
	LightningActivityLevel           GridpointForecastTimeSeries `json:"lightningActivityLevel"`
	TwentyFootWindSpeed              GridpointForecastTimeSeries `json:"twentyFootWindSpeed"`
	TwentyFootWindDirection          GridpointForecastTimeSeries `json:"twentyFootWindDirection"`
	WaveHeight                       GridpointForecastTimeSeries `json:"waveHeight"`
	WavePeriod                       GridpointForecastTimeSeries `json:"wavePeriod"`
	WaveDirection                    GridpointForecastTimeSeries `json:"waveDirection"`
	PrimarySwellHeight               GridpointForecastTimeSeries `json:"primarySwellHeight"`
	PrimarySwellDirection            GridpointForecastTimeSeries `json:"primarySwellDirection"`
	SecondarySwellHeight             GridpointForecastTimeSeries `json:"secondarySwellHeight"`
	SecondarySwellDirection          GridpointForecastTimeSeries `json:"secondarySwellDirection"`
	WavePeriod2                      GridpointForecastTimeSeries `json:"wavePeriod2"`
	WindWaveHeight                   GridpointForecastTimeSeries `json:"windWaveHeight"`
	DispersionIndex                  GridpointForecastTimeSeries `json:"dispersionIndex"`
	Pressure                         GridpointForecastTimeSeries `json:"pressure"`
	ProbabilityOfTropicalStormWinds  GridpointForecastTimeSeries `json:"probabilityOfTropicalStormWinds"`
	ProbabilityOfHurricaneWinds      GridpointForecastTimeSeries `json:"probabilityOfHurricaneWinds"`
	PotentialOf15mphWinds            GridpointForecastTimeSeries `json:"potentialOf15mphWinds"`
	PotentialOf25mphWinds            GridpointForecastTimeSeries `json:"potentialOf25mphWinds"`
	PotentialOf35mphWinds            GridpointForecastTimeSeries `json:"potentialOf35mphWinds"`
	PotentialOf45mphWinds            GridpointForecastTimeSeries `json:"potentialOf45mphWinds"`
	PotentialOf20mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf20mphWindGusts"`
	PotentialOf30mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf30mphWindGusts"`
	PotentialOf40mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf40mphWindGusts"`
	PotentialOf50mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf50mphWindGusts"`
	PotentialOf60mphWindGusts        GridpointForecastTimeSeries `json:"potentialOf60mphWindGusts"`
	GrasslandFireDangerIndex         GridpointForecastTimeSeries `json:"grasslandFireDangerIndex"`
	ProbabilityOfThunder             GridpointForecastTimeSeries `json:"probabilityOfThunder"`
	DavisStabilityIndex              GridpointForecastTimeSeries `json:"davisStabilityIndex"`
	AtmosphericDispersionIndex       GridpointForecastTimeSeries `json:"atmosphericDispersionIndex"`
	LowVisibilityOccurrenceRiskIndex GridpointForecastTimeSeries `json:"lowVisibilityOccurrenceRiskIndex"`
	Stability                        GridpointForecastTimeSeries `json:"stability"`
	RedFlagThreatIndex               GridpointForecastTimeSeries `json:"redFlagThreatIndex"`
	Point                            *PointsResponse
}

// GridpointForecastTimeSeriesValue holds the JSON value for a
// GridpointForecastTimeSeries' values[x] item.
type GridpointForecastTimeSeriesValue struct {
	ValidTime string  `json:"validTime"` // ISO 8601 time interval, e.g. 2019-07-04T18:00:00+00:00/PT3H
	Value     float64 `json:"value"`
}

// GridpointForecastTimeSeries holds a series of data from a gridpoint forecast
type GridpointForecastTimeSeries struct {
	Uom    string                             `json:"uom"` // Unit of Measure
	Values []GridpointForecastTimeSeriesValue `json:"values"`
}

type Observation struct {
	ID              string `json:"@id"`
	IDType          string `json:"@type"`
	Icon            string `json:"icon"`
	TextDescription string
	Station         string
	Elevation       QuantitativeValue
	Geometry        string // ex: POINT(-73.09 40.79)

	Temperature               QuantitativeValue
	MaxTemperatureLast24Hours QuantitativeValue
	MinTemperatureLast24Hours QuantitativeValue

	PrecipitationLastHour   QuantitativeValue
	PrecipitationLast3Hours QuantitativeValue
	PrecipitationLast6Hours QuantitativeValue

	RelativeHumidity   QuantitativeValue
	SeaLevelPressure   QuantitativeValue
	Dewpoint           QuantitativeValue
	BarometricPressure QuantitativeValue

	WindChill     QuantitativeValue
	WindGust      QuantitativeValue
	WindDirection QuantitativeValue
	WindSpeed     QuantitativeValue

	HeatIndex  QuantitativeValue
	Visibility QuantitativeValue

	RawMessage string
	Timestamp  time.Time

	CloudLayers    []CloudLayerReading
	PresentWeather []interface{}
}

type CloudLayerReading struct {
	Amount string
	Base   QuantitativeValue
}

type ObservationsResponse struct {
	Observations []Observation `json:"@graph"`
}
