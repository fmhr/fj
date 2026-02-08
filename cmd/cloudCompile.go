package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/fmhr/fj/cmd/setup"
)

func checkConfigCloudCompile(config *setup.Config) error {
	if config.Language == "" {
		return fmt.Errorf("error: [Language] must not be empty")
	}
	if config.SourceFilePath == "" {
		return fmt.Errorf("error: [SourcePath] must not be empty")
	}
	if config.BinaryPath == "" {
		return fmt.Errorf("error: [BinaryPath] must not be empty")
	}
	if config.CompilerURL == "" {
		return fmt.Errorf("error: [CompilerURL] must not be empty")
	}
	return nil
}

// コンパイルに必要なもの
//  1. ソースファイル受け取る
//  2. ソースファイルの名前
//  3. コンパイルコマンド
//  4. file=ソースファイル, compileCmd=コンパイルコマンド,
//     srcFile=ソースファイル名, binaryFile=バイナリファイル名
func CloudCompile(config *setup.Config) error {
	_, ok := setup.LanguageSets[config.Language]
	if !ok {
		return fmt.Errorf("error: language [%s] is not supported. suported %v", config.Language, setup.LanguageList())
	}
	log.Println("cloud compiling...")
	if err := checkConfigCloudCompile(config); err != nil {
		return err
	}
	// ソースファイルを開く
	file, err := os.Open(config.SourceFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// マルチパートフォームを作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(config.SourceFilePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.WriteField("language", config.Language)
	writer.WriteField("sourcePath", config.SourceFilePath)
	writer.WriteField("binaryPath", config.BinaryPath)
	writer.WriteField("bucket", config.Bucket)

	writer.Close()

	// POSTリクエストを送信
	req, err := http.NewRequest("POST", config.CompilerURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("error response status code:%d resp:%s", resp.StatusCode, string(bodyBytes))
	}

	// cloud storageに保存したバイナルの名前を取得
	content := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(content)
	if err != nil {
		return err
	}
	filename, err := url.QueryUnescape(params["filename"])
	if err != nil {
		log.Println("error: failed to unescape filename")
		return err
	}
	config.TmpBinary = filename
	log.Println("cloud compile done: bucket:", filename)
	return nil
}
