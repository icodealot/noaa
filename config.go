package noaa

import (
	"fmt"
	"strings"
)

// deprecated
// Default values for the weather.gov REST API config which will
// be replaced by Config. These are subject to deletion in the future.
// Instead, use noaa.GetConfig followed by:
//     Config.BaseURL, Config.UserAgent, Config.Accept
const (
	API       = "https://api.weather.gov"
	APIKey    = "github.com/icodealot/noaa" // User-Agent default value
	APIAccept = "application/ld+json"       // Changes may affect struct mappings below
)

const (
	templateEndpointOffices = "%s/offices/%s"   // base url, office id
	templateEndpointPoints  = "%s/points/%s,%s" // base url, lat, lon
)

// Config instance for the API calls executed by the NOAA client.
var config = GetDefaultConfig()

// Config describes important values for the NOAA API and allows for
// configuration and testing of various options. Note, the User-Agent
// field of HTTP requests serves as a proxy for an API key and in the
// future weather.gov might change this behavior.
// See http://www.weather.gov/documentation/services-web-api
type Config struct {
	BaseURL   string `json:"baseUrl"` // Do not include a trailing slash
	UserAgent string `json:"apiKey"`  // ex. (myweatherapp.com, contact@myweatherapp.com)
	Accept    string `json:"accept"`  // application/geo+json, etc. defaults to ld+json
	Units     string `json:"units"`   // "us" (the default if blank) or "si" for metric
}

func (c *Config) endpointOffices(id string) string {
	return fmt.Sprintf(templateEndpointOffices, config.BaseURL, id)
}

func (c *Config) endpointPoints(lat string, lon string) string {
	return fmt.Sprintf(templateEndpointPoints, config.BaseURL, lat, lon)
}

func (c *Config) getUnitsQueryParam(prefix string) string {
	queryParam := ""
	if config.Units != "" {
		queryParam = prefix + "units=" + config.Units
	}
	return queryParam
}

// SetUserAgent changes the string used for the User-Agent header when making
// requests. See https://www.weather.gov/documentation/services-web-api
// (Authentication) for details.  By default, this module uses a github.com URL.
func SetUserAgent(userAgent string) {
	if len(userAgent) == 0 {
		panic("the api requires a user-agent")
	}
	config.UserAgent = userAgent
}

// SetUnits can be used to change the units returned by the weather.gov API from
// US to metric. By default, if no units are specified, then the API assumes US.
func SetUnits(uom string) {
	units := strings.ToLower(uom)
	if units != "us" && units != "si" {
		config.Units = ""
	} else {
		config.Units = units
	}
}

// SetConfig replaces the config with all new values in one call. The individual
// Set* functions can also be used to replace only specified values.
func SetConfig(c Config) {
	if !isConfigValid(c) {
		panic("invalid configuration")
	}
	config = c
}

// GetConfig is used to return the current configuration of the client. This allows
// for testing and inspection as needed.
func GetConfig() Config {
	return config
}

// GetDefaultConfig returns a config struct that can be used as a starting point
// for configuration changes. See examples in `example_test.go`.
func GetDefaultConfig() Config {
	return Config{
		BaseURL:   API,
		UserAgent: APIKey,
		Accept:    APIAccept,
		Units:     "", // defaults to US units if unspecified
	}
}

// SetBaseURL changes the base URL of the API. This can be useful for testing
// and if the weather.gov endpoint is relocated, in a pinch you could set it.
// Probably not useful in general.
func SetBaseURL(url string) {
	if len(url) == 0 {
		panic("the api requires a base url")
	}
	config.BaseURL = url
}

// SetAcceptHeader changes the format of the response. Note, this is largely a
// placeholder for future use and testing as the Go types defined in this wrapper
// assume application/ld+json. Using anything else is undefined.
// Probably not useful in general.
func SetAcceptHeader(accept string) {
	if len(accept) == 0 {
		panic("the api requires an accept header")
	}
	config.Accept = accept
}

// isConfigValid determines whether the provided config might be valid. Under
// certain conditions we can determine if the config is definitely not valid.
func isConfigValid(c Config) bool {
	if len(c.Units) > 0 && c.Units != "us" && c.Units != "si" {
		return false
	}
	if len(c.Accept) == 0 || len(c.BaseURL) == 0 || len(c.UserAgent) == 0 {
		return false
	}
	return true
}
