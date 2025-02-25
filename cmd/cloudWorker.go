package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd/setup"
)

// requestToWorker はバイナリをワーカーに送信する
func requestToWorker(config *setup.Config, seed int) (*orderedmap.OrderedMap[string, string], error) {
	start := time.Now()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// configをJSONに変換
	configData, err := json.Marshal(config)
	if err != nil {
		return nil, NewStackTraceError(err.Error())
	}
	// JSON configを追加
	configPart, err := writer.CreateFormField("config")
	if err != nil {
		return nil, NewStackTraceError(err.Error())
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
			return nil, WrapError(err)
		}
		return nil, NewStackTraceError(fmt.Sprintf("error response status code:%d resp:%s", resp.StatusCode, string(bodyBytes)))
	}

	// レスポンスボディから文字列を取り出す
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	rtn := orderedmap.NewOrderedMap[string, string]()
	var tempMap map[string]string
	if err := json.Unmarshal(bodyBytes, &tempMap); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}
	for key, value := range tempMap {
		rtn.Set(key, value)
	}
	elapsed := time.Since(start)
	rtn.Set("responseTime", elapsed.String())
	score, ok := rtn.Get("Score")
	if !ok {
		return nil, fmt.Errorf("failed to get score from response body: %v", err)
	}
	if score == "0" {
		log.Println("Score=0:response body:", string(bodyBytes))
	}
	return rtn, nil
}

func SendBinaryToWorker(config *setup.Config, seed int, binaryNameInBucket string) (*orderedmap.OrderedMap[string, string], error) {
	if config.WorkerURL == "" {
		return nil, NewStackTraceError("worker URL is not specified")
	}
	if config.BinaryPath == "" {
		return nil, NewStackTraceError("binary path is not specified")
	}
	return requestToWorker(config, seed)
}
