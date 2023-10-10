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

// コンパイルに必要なもの
//  1. ソースファイル受け取る
//  2. ソースファイルの名前
//  3. コンパイルコマンド
//  4. file=ソースファイル, compileCmd=コンパイルコマンド,
//     srcFile=ソースファイル名, binaryFile=バイナリファイル名
func CloudCompile(config *Config) (string, error) {
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
	writer.CreateFormFile("srcFile", file.Name())
	writer.CreateFormFile("compileCmd", config.CompileCmd)
	writer.CreateFormFile("binaryFile", binaryName)

	writer.Close()

	// POSTリクエストを送信
	req, err := http.NewRequest("POST", config.CompilerURL, body)
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
