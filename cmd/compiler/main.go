package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/fmhr/fj"
)

// コンパイルに必要なもの
//  1. ソースファイル受け取る
//  2. ソースファイルの名前
//  3. コンパイルコマンド
//  4. file=ソースファイル, compileCmd=コンパイルコマンド,
//     srcFile=ソースファイル名, binaryFile=バイナリファイル名
func compileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// マルチパートリーダー
	err := r.ParseMultipartForm(10 << 20) //10MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	// ソースファイル
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// コンパイルコマンド
	language := r.FormValue("language")
	if language == "" {
		http.Error(w, "Language not specified", http.StatusBadRequest)
		return
	}
	// ソースファイル名
	source := r.FormValue("sourcePath")
	if source == "" {
		http.Error(w, "Source file not specified", http.StatusBadRequest)
		return
	}
	err = createFileWithDirs(source, nil)
	if err != nil {
		http.Error(w, "Failed to create source file", http.StatusInternalServerError)
		return
	}

	// バイナリファイル名
	binaryFileName := r.FormValue("binaryPath")
	err = createFileWithDirs(binaryFileName, nil)
	if err != nil {
		http.Error(w, "Failed to create binary file", http.StatusInternalServerError)
		return
	}

	// ソースファイルをディスクに書き込む
	srcBytes, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("Failed to read the uploaded file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	err = os.WriteFile(source, srcBytes, 0644)
	if err != nil {
		msg := fmt.Sprintf("Failed to write the uploaded file to disk: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	compileCmd := fj.LanguageSets[language].CompileCmd
	// コンパイル
	cmds := strings.Fields(compileCmd)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	msg, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
		msg := fmt.Sprintf("Failed to compile: [%s]%v msg: %s", cmd.String(), err, string(msg))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	bucketName := r.FormValue("bucket")

	// google cloud storageにバイナリとソースファイルをアップロード
	rdm := generateRandomString(10)
	newfilename, err := uploarFileToGoogleCloudStorage(bucketName, binaryFileName, rdm)
	if err != nil {
		http.Error(w, "Failed to upload binary file to bucket", http.StatusInternalServerError)
		return
	}
	_, err = uploarFileToGoogleCloudStorage(bucketName, source, rdm)
	if err != nil {
		http.Error(w, "Failed to upload source file to bucket", http.StatusInternalServerError)
		return
	}

	// バイナリの読み込み
	binary, err := os.ReadFile(binaryFileName)
	if err != nil {
		http.Error(w, "Failed to read compiled binary", http.StatusInternalServerError)
		return
	}
	// ソースファイルとバイナリファイルを削除
	defer os.Remove(source)
	defer os.Remove(binaryFileName)

	// バイナリをレスポンスとして返す
	w.Header().Set("Content-Disposition", "attachment; filename="+newfilename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(binary)
	if err != nil {
		http.Error(w, "Failed to write binary to response", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/compiler", compileHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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

func uploarFileToGoogleCloudStorage(bucketName string, file, rdm string) (string, error) {
	// ファイルの読み込み
	f, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	newfilename := rdm + "-" + file

	wc := client.Bucket(bucketName).Object(newfilename).NewWriter(ctx)
	_, err = wc.Write(f)
	if err != nil {
		return "", fmt.Errorf("failed to write object: %v", err)
	}
	defer wc.Close()

	return newfilename, nil
}

func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		log.Println(err)
	}
	for i := 0; i < length; i++ {
		b[i] = letters[b[i]%byte(len(letters))]
	}
	return string(b)
}
