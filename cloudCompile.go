package fj

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func checkConfigCloudCompile(config *Config) error {
	if config.Source == "" {
		return ErrorTrace("error: [SourcePath] must not be empty", nil)
	}
	if config.Binary == "" {
		return ErrorTrace("error: [BinaryPath] must not be empty", nil)
	}
	if config.CompileCmd == "" {
		return ErrorTrace("error: [CompileCmd] must not be empty", nil)
	}
	if config.CompilerURL == "" {
		return ErrorTrace("error: [CompilerURL] must not be empty", nil)
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
	log.Println("cloud compiling...")
	if err := checkConfigCloudCompile(config); err != nil {
		return "", ErrorTrace("invalid config: %w", err)
	}
	// ソースファイルを開く
	file, err := os.Open(config.Source)
	if err != nil {
		return "", ErrorTrace("error opening file: %w", err)
	}
	defer file.Close()
	// マルチパートフォームを作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(config.Source))
	if err != nil {
		return "", ErrorTrace("error creating form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", ErrorTrace("error writing to form file: %w", err)
	}
	writer.WriteField("sourcePath", config.Source)
	writer.WriteField("compileCmd", config.CompileCmd)
	writer.WriteField("binaryPath", config.Binary)

	writer.Close()

	// POSTリクエストを送信
	req, err := http.NewRequest("POST", config.CompilerURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return "", ErrorTrace("error making new requets: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", ErrorTrace("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", ErrorTrace("error reading response body: %w", err)
		}
		fmt.Print("server response:", string(bodyBytes))
		return "", ErrorTrace(string(bodyBytes), fmt.Errorf("error response status code: %d", resp.StatusCode))
	}

	// 受信後
	// バイナリを保存 保存場所はOSの一時フォルダ
	out, err := os.CreateTemp("", "binary-*")
	if err != nil {
		return "", ErrorTrace("error createing output file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", ErrorTrace("error saving binary: %w", err)
	}
	log.Println("binary saved to", out.Name())
	return out.Name(), nil
}
