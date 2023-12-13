package fj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// requestToWorker はバイナリをワーカーに送信する
func requestToWorker(config *Config, seed int) (rtn map[string]float64, err error) {
	start := time.Now()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// configをJSONに変換
	configData, err := json.Marshal(config)
	if err != nil {
		return nil, ErrorTrace("failed to marshal config: %v", err)
	}
	// JSON configを追加
	configPart, err := writer.CreateFormField("config")
	if err != nil {
		return nil, ErrorTrace("failed to create form field for config: %v", err)
	}
	configPart.Write(configData)

	writer.WriteField("seed", fmt.Sprintf("%d", seed))
	writer.Close()

	// リクエストの作成
	req, err := http.NewRequest("POST", config.WorkerURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// リクエストの送受信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send  HTTP request to worker: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrorTrace("error reading response body: %w", err)
		}
		return nil, ErrorTrace(fmt.Sprintf("error response status code:%d resp:%s", resp.StatusCode, string(bodyBytes)), err)
	}

	// レスポンスボディから文字列を取り出す
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	if err := json.Unmarshal(bodyBytes, &rtn); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}
	elapsed := time.Since(start)
	rtn["responseTime"] = elapsed.Seconds()
	return rtn, nil
}

func SendBinaryToWorker(config *Config, seed int, binaryNameInBucket string) (rtn map[string]float64, err error) {
	if config.WorkerURL == "" {
		return nil, ErrorTrace("", fmt.Errorf("worker URL is not specified"))
	}
	if config.Binary == "" {
		return nil, ErrorTrace("", fmt.Errorf("binary path is not specified"))
	}
	return requestToWorker(config, seed)
}
