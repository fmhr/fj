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

	// バイナリをCloud Storageからダウンロード
	err = downloadFileFromGoogleCloudStorage(config.Bucket, config.tmpBinary, config.Binary)
	if err != nil {
		errmsg := fmt.Sprint("Failed to download binary from Cloud Storage:", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	// バイナリの受け取り
	file, _, err := r.FormFile("binary")
	if err != nil {
		errmsg := fmt.Sprint("Failed to get the binary:", err.Error())
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}
	defer file.Close()

	if config.Binary == "" {
		http.Error(w, "BinaryName is empty", http.StatusInternalServerError)
		return
	}
	binaryPath, err := os.Create(config.Binary)
	if err != nil {
		errmsg := fmt.Sprintf("Failed to Create the binary file: %v", err)
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(binaryPath, file)
	if err != nil {
		errmsg := fmt.Sprint("Failed to copy the binary to the temp file:", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	// 実行権限を与える
	err = os.Chmod(binaryPath.Name(), 0755)
	if err != nil {
		errmsg := fmt.Sprint("Failed to chmod", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	binaryPath.Close()

	// seed
	seedString := r.FormValue("seed")
	seedInt, err := strconv.Atoi(seedString)
	if err != nil {
		errmsg := fmt.Sprint("Failed to convert seed to int", err.Error())
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}
	out, err := exexute(&config, seedInt)
	if err != nil {
		errmsg := fmt.Sprint("Failed to execute", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	// json
	jsonData, err := json.Marshal(out)
	if err != nil {
		errmsg := fmt.Sprint("Failed to marshal", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/worker", handler)
	http.ListenAndServe(":8080", nil)
}

func createFileWithDirs(path string, data []byte) error {
	dir := filepath.Dir(path)

	// ディレクトリが存在しない場合は作成
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// ファイルを作成
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	return nil
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
