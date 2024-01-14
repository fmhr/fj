package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/fmhr/fj"
)

func main() {
	http.HandleFunc("/worker", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	// マルチパートリーダー
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		errmsg := fmt.Sprint("Failed to parse multipart form:", err.Error())
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}

	// Configの受け取り
	configPart := r.FormValue("config")
	var config fj.Config
	err = json.Unmarshal([]byte(configPart), &config)
	if err != nil {
		errmsg := fmt.Sprint("Failed to unmarshal config:", err.Error())
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}

	// すでにバイナリがあるか確認
	if filepath.Clean(config.TmpBinary) != config.TmpBinary {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}
	_, err = os.Stat(filepath.Clean(config.TmpBinary))

	if os.IsNotExist(err) {
		// バイナリをCloud Storageからダウンロード
		if config.Bucket == "" {
			http.Error(w, "BucketName is empty", http.StatusInternalServerError)
			return
		}
		err = downloadFileFromGoogleCloudStorage(config.Bucket, config.TmpBinary, config.Binary)
		if err != nil {
			errmsg := fmt.Sprint("Failed to download binary from Cloud Storage:", err.Error())
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}

		// 実行権限を与える
		err = os.Chmod(config.Binary, 0755)
		if err != nil {
			errmsg := fmt.Sprint("Failed to chmod", err.Error())
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}
	} else {
		// tmpバイナリをmainに改名
		err = os.Rename(config.TmpBinary, config.Binary)
		if err != nil {
			errmsg := fmt.Sprint("Failed to rename binary", err.Error())
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}
	}

	// 入力ファイルを作成
	seedString := r.FormValue("seed")
	seedInt, err := strconv.Atoi(seedString)
	if err != nil {
		errmsg := fmt.Sprint("Failed to convert seed to int", err.Error())
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}
	// 実行
	out, err := exexute(&config, seedInt)
	if err != nil {
		errmsg := fmt.Sprint("Failed to execute", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	// json
	jsonData, err := json.Marshal((*fj.EncodableOrderedMap)(out))
	if err != nil {
		errmsg := fmt.Sprint("Failed to marshal", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	// バイナリをtmpネームに改名
	err = os.Rename(config.Binary, config.TmpBinary)
	if err != nil {
		errmsg := fmt.Sprint("Failed to rename binary", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
}

func downloadFileFromGoogleCloudStorage(bucketName string, objectName string, destination string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create reader: %v", err)
	}
	defer rc.Close()

	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, rc); err != nil {
		return fmt.Errorf("failed to copy: %v", err)
	}

	return nil
}
