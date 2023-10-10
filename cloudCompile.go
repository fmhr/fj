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

func checkConfigCloudCompile(config *Config) error {
	if config.SourcePath == "" {
		return fmt.Errorf("error: [SourcePath] must not be empty")
	}
	if config.BinaryPath == "" {
		return fmt.Errorf("error: [BinaryPath] must not be empty")
	}
	if config.CompileCmd == "" {
		return fmt.Errorf("error: [CompileCmd] must not be empty")
	}
	if config.CompilerURL == "" {
		return TraceError("error: [CompilerURL] must not be empty")
	}
	return nil
}

// コンパイルに必要なもの
//  1. ソースファイル受け取る
//  2. ソースファイルの名前
//  3. コンパイルコマンド
//  4. file=ソースファイル, compileCmd=コンパイルコマンド,
//     srcFile=ソースファイル名, binaryFile=バイナリファイル名
func CloudCompile(config *Config) (string, error) {
	if err := checkConfigCloudCompile(config); err != nil {
		return "", fmt.Errorf("invalid config: %w", err)
	}
	// ソースファイルを開く
	file, err := os.Open(config.SourcePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	// マルチパートフォームを作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(config.SourcePath))
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error writing to form file: %w", err)
	}
	// filename
	binaryName := filepath.Base(config.BinaryPath)
	writer.WriteField("srcFile", file.Name())
	writer.WriteField("compileCmd", config.CompileCmd)
	writer.WriteField("binaryFile", binaryName)

	writer.Close()

	// POSTリクエストを送信
	req, err := http.NewRequest("POST", config.CompilerURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return "", TraceErrorf("error making new requets: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", TraceErrorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", TraceErrorf("error reading response body: %w", err)
		}
		fmt.Print("server response:", string(bodyBytes))
		return "", TraceError(fmt.Sprintf("error response status code: %d", resp.StatusCode))
	}

	// 受信後
	out, err := os.CreateTemp("", "binary-*")
	if err != nil {
		return "", fmt.Errorf("error createing output file: %w", err)
	}
	defer out.Close()

	// バイナリを保存
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error saving binary: %w", err)
	}

	return out.Name(), nil
}
