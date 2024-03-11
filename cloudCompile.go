package fj

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func checkConfigCloudCompile(config *Config) error {
	if config.Source == "" {
		return NewStackTraceError("error: [SourcePath] must not be empty")
	}
	if config.Binary == "" {
		return NewStackTraceError("error: [BinaryPath] must not be empty")
	}
	if config.CompileCmd == "" {
		return NewStackTraceError("error: [CompileCmd] must not be empty")
	}
	if config.CompilerURL == "" {
		return NewStackTraceError("error: [CompilerURL] must not be empty")
	}
	return nil
}

// コンパイルに必要なもの
//  1. ソースファイル受け取る
//  2. ソースファイルの名前
//  3. コンパイルコマンド
//  4. file=ソースファイル, compileCmd=コンパイルコマンド,
//     srcFile=ソースファイル名, binaryFile=バイナリファイル名
func CloudCompile(config *Config) error {
	_, ok := LanguageSets[config.Language]
	if !ok {
		return NewStackTraceError(fmt.Sprintf("error: language [%s] is not supported. suported %v", config.Language, languageList()))
	}
	log.Println("cloud compiling...")
	if err := checkConfigCloudCompile(config); err != nil {
		return err
	}
	// ソースファイルを開く
	file, err := os.Open(config.Source)
	if err != nil {
		return err
	}
	defer file.Close()
	// マルチパートフォームを作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(config.Source))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.WriteField("language", config.Language)
	writer.WriteField("sourcePath", config.Source)
	writer.WriteField("compileCmd", config.CompileCmd)
	writer.WriteField("binaryPath", config.Binary)
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
		msg := fmt.Sprint("server response:", string(bodyBytes), "url:", config.CompilerURL, "\n")
		return NewStackTraceError(msg)
	}

	// cloud storageに保存したバイナルの名前を取得
	content := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(content)
	if err != nil {
		return err
	}
	filename := params["filename"]
	config.TmpBinary = filename
	log.Println("cloud compile done:", filename)
	return nil
}
