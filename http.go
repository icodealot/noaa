package noaa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// HTTP GET the noaa endpoint provided. We could just use http.Get() but
// this helps since we include some custom header values
func get(endpoint string) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", config.Accept)
	req.Header.Add("User-Agent", config.UserAgent)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%d %s", res.StatusCode, res.Status))
	}
	return res, nil
}

func decode(endpoint string, v any) error {
	res, err := get(endpoint)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(v); err != nil {
		return err
	}
	return nil
}
