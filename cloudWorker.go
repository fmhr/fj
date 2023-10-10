package fj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func sendBinaryToWorker(workerURL, binaryPath, language string, seed int) (rtn map[string]float64, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// バイナリの追加
	file, err := os.Open(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open binary file %s: %v", binaryPath, err)
	}
	defer file.Close()
	part, err := writer.CreateFormFile("binary", filepath.Base(binaryPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file for binary: %v", err)
	}
	// バイナリの書き込み
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to write binary to form file: %v", err)
	}
	// filename
	writer.CreateFormFile("filename", file.Name())
	// language
	writer.WriteField("langiage", language)
	// seed
	writer.WriteField("seed", fmt.Sprintf("%d", seed))

	writer.Close()

	// リクエストの送信
	req, err := http.NewRequest("POST", workerURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send  HTTP request to worker: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("worker returned an unexpected status: %v", resp.Status)
	}

	// レスポンスボディから文字列を取り出す
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	if err := json.Unmarshal(bodyBytes, &rtn); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}

	return rtn, nil
}
