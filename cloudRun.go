package fj

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CloudRun(cnf *Config, seed int) (map[string]float64, error) {
	return googleCloudRun(cnf, seed)
}

func googleCloudRun(cnf *Config, seed int) (map[string]float64, error) {
	url := fmt.Sprintf("%s?seed=%d", cnf.CloudURL, seed)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var data map[string]float64
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error parsing response body: %v", err)
	}

	return data, nil
}
