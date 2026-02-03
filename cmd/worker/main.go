package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/fmhr/fj/cmd"
	"github.com/fmhr/fj/cmd/setup"
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
	var config setup.Config
	err = json.Unmarshal([]byte(configPart), &config)
	if err != nil {
		errmsg := fmt.Sprint("Failed to unmarshal config:", err.Error())
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}

	// javaやpythonではソースフィルをダウンロードする
	if config.Language == "java" || config.Language == "C#" {
		err = downloadFileFromGoogleCloudStorage(config.Bucket, config.TmpBinary, config.SourceFilePath)
		if err != nil {
			errmsg := fmt.Sprint("Failed to download source file from Cloud Storage:", err.Error())
			errmsg += fmt.Sprint("Bucket:", config.Bucket, " TmpBinary:", config.TmpBinary, " SourceFilePath:", config.SourceFilePath)
			http.Error(w, errmsg, http.StatusInternalServerError)
			return
		}
		// javaの場合はコンパイル
		if config.Language == "java" || config.Language == "C#" {
			compileCmd := setup.LanguageSets[config.Language].CompileCmd
			cmds := strings.Fields(compileCmd)
			cmd := exec.Command(cmds[0], cmds[1:]...)
			msg, err := cmd.CombinedOutput()
			if err != nil {
				log.Println(msg)
				err := fmt.Errorf("failed to compile: [%s]%v msg: %s", cmd.String(), err, string(msg))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if config.Language == "C#" {
				// 実行権限を与える
				err = os.Chmod(config.BinaryPath, 0755)
				if err != nil {
					errmsg := fmt.Sprint("Failed to chmod:", err.Error())
					http.Error(w, errmsg, http.StatusInternalServerError)
					return
				}
			}
		}
	} else {
		// すでにバイナリがあるか確認
		if filepath.Clean(config.TmpBinary) != config.TmpBinary {
			http.Error(w, "Invalid file path:", http.StatusBadRequest)
			return
		}
		tmpBinaryFileName := filepath.Clean(config.TmpBinary)
		_, err = os.Stat(tmpBinaryFileName)

		if os.IsNotExist(err) {
			// バイナリをCloud Storageからダウンロード
			err = downloadFileFromGoogleCloudStorage(config.Bucket, config.TmpBinary, config.BinaryPath)
			if err != nil {
				errmsg := fmt.Sprint("Failed to download binary from Cloud Storage:", err.Error())
				errmsg += fmt.Sprint("Bucket:", config.Bucket, " TmpBinary:", config.TmpBinary, " BinaryPath:", config.BinaryPath)
				http.Error(w, errmsg, http.StatusInternalServerError)
				return
			}

			// 実行権限を与える
			err = os.Chmod(config.BinaryPath, 0755)
			if err != nil {
				errmsg := fmt.Sprint("Failed to chmod", err.Error())
				http.Error(w, errmsg, http.StatusInternalServerError)
				return
			}
		} else {
			// tmpバイナリをmainに改名
			err = os.Rename(config.TmpBinary, config.BinaryPath)
			if err != nil {
				errmsg := fmt.Sprint("Failed to rename binary", err.Error())
				http.Error(w, errmsg, http.StatusInternalServerError)
				return
			}
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
	out, err := execute(&config, seedInt)
	if err != nil {
		errmsg := fmt.Sprint("Failed to execute", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	// json
	jsonData, err := json.Marshal((*cmd.EncodableOrderedMap)(out))
	if err != nil {
		errmsg := fmt.Sprint("Failed to marshal", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	// バイナリをtmpネームに改名
	// C#はディレクトリ以下にバイナリがあるので、ディレクトリごとつくる
	if err := os.MkdirAll(filepath.Dir(config.TmpBinary), 0755); err != nil {
		log.Println(err)
	}
	err = os.Rename(config.BinaryPath, config.TmpBinary)
	if err != nil {
		errmsg := fmt.Sprint("Failed to rename binary", err.Error())
		log.Println(errmsg, config.BinaryPath, config.TmpBinary)
		//http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
}

func downloadFileFromGoogleCloudStorage(bucketName string, objectName string, destination string) error {
	// ローカルパスのディレクトリを確認し、存在しない場合は作成
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return err
	}
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
