package fj

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func uploadAndReceive(serverURL, sourcFile, language string) (string, error) {
	file, err := os.Open(sourcFile)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	// マルチパートフォームを作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(sourcFile))
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error writing to form file: %w", err)
	}
	writer.WriteField("language", language)
	writer.Close()

	// POSTリクエストを送信
	req, err := http.NewRequest("POST", serverURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return "", fmt.Errorf("error making new requets: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned an error: %s", resp.Status)
	}

	// バイナリを一時ディレクトリに保存
	tempDir, err := os.MkdirTemp("", "binary-*")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %w", err)
	}
	outputPath := filepath.Join(tempDir, "a.out")

	out, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("error createing output file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error saving binary: %w", err)
	}

	return outputPath, nil
}
