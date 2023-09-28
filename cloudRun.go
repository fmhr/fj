package fj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pelletier/go-toml/v2"
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

	conData, err := toml.Marshal(cnf)
	if err != nil {
		return nil, fmt.Errorf("error marshaling config: %v", err)
	}

	resp, err := http.Post(finalURL, "application/toml", bytes.NewReader(conData))
	if err != nil {
		return nil, fmt.Errorf("error making POST request to %s: %v", finalURL, err)
	}
	defer resp.Body.Close()

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
	//fmt.Println(mapString(data))

	return data, nil
}
