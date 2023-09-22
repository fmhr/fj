package fj

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func CloudRun(cnf *Config, seed int) (map[string]float64, error) {
	return cloudRun(cnf, seed)
}

func cloudRun(cnf *Config, seed int) (map[string]float64, error) {
	baseURL, err := url.Parse(cnf.CloudURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing cloud url: %v", err)
	}
	params := url.Values{}
	params.Add("seed", strconv.Itoa(seed))

	baseURL.RawQuery = params.Encode()
	finalURL := baseURL.String()

	resp, err := http.Get(finalURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %v", finalURL, err)
	}
	defer resp.Body.Close()

	log.Println(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error making GET request to %s: %s", finalURL, resp.Status)
	}

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
