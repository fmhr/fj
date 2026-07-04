package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/fmhr/fj/cmd/setup"
)

// コンパイルに必要なもの
//  1. ソースファイル受け取る
//  2. ソースファイルの名前
//  3. コンパイルコマンド
//  4. file=ソースファイル, compileCmd=コンパイルコマンド,
//     srcFile=ソースファイル名, binaryFile=バイナリファイル名
func parseCompileRequest(r *http.Request) (file io.ReadCloser, language string, bucket string, err error) {
	if r.Method != http.MethodPost {
		return nil, "", "", fmt.Errorf("Only POST method is allowed")
	}
	// マルチパートリーダー
	err = r.ParseMultipartForm(50 << 20) //10MB limit
	if err != nil {
		return nil, "", "", fmt.Errorf("Failed to parse form:%w", err)
	}

	// ソースファイル
	file, _, err = r.FormFile("file")
	if err != nil {
		return nil, "", "", fmt.Errorf("Failed to get the file:%w", err)
	}

	// language
	language = r.FormValue("language")
	if language == "" {
		file.Close()
		return nil, "", "", fmt.Errorf("Language not specified:")
	}

	if _, ok := setup.LanguageSets[language]; !ok {
		file.Close()
		return nil, "", "", fmt.Errorf("Language [%s] is not supported.", language)
	}

	bucket = r.FormValue("bucket")
	return file, language, bucket, nil
}

func prepareAndSaveSource(file io.Reader, language string) (source string, binaryFileName string, err error) {
	// ソースファイル名
	source = setup.LanguageSets[language].FileName
	if source == "" {
		return "", "", fmt.Errorf("Source file not specified:")
	}
	err = createFileWithDirs(source, nil)
	if err != nil {
		return "", "", fmt.Errorf("Failed to create source file:%w", err)
	}

	// バイナリファイル名
	binaryFileName = setup.LanguageSets[language].BinaryPath
	err = createFileWithDirs(binaryFileName, nil)
	if err != nil {
		return "", "", fmt.Errorf("Failed to create binary file:%w", err)
	}

	// ソースファイルをディスクに書き込む
	srcBytes, err := io.ReadAll(file)
	if err != nil {
		return "", "", fmt.Errorf("Failed to read the uploaded file: %w", err)
	}
	err = os.WriteFile(source, srcBytes, 0644)
	if err != nil {
		return "", "", fmt.Errorf("Failed to write the uploaded file to disk: %w", err)
	}

	return source, binaryFileName, nil
}

func runCompiler(language string) ([]byte, error) {
	compileCmd := setup.LanguageSets[language].CompileCmd
	cmds := strings.Fields(compileCmd)
	if len(cmds) == 0 {
		return nil, fmt.Errorf("Compile command is empty")
	}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	msg, err := cmd.CombinedOutput()
	if err != nil {
		return msg, fmt.Errorf("Failed to compile: [%s]%w", cmd.String(), err)
	}
	return msg, nil
}

func uploadArtifacts(bucketName string, source string, binaryFileName string, rdm string) (string, error) {
	// バイナリのアップロード
	newfilename, err := uploadFileToGoogleCloudStorage(bucketName, binaryFileName, rdm)
	if err != nil {
		return "", fmt.Errorf("Failed to upload binary file to bucket:%w", err)
	}
	// ソースコードのアップロード
	_, err = uploadFileToGoogleCloudStorage(bucketName, source, rdm)
	if err != nil {
		return "", fmt.Errorf("Failed to upload source file to bucket:%w", err)
	}
	return newfilename, nil
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	file, language, bucketName, err := parseCompileRequest(r)
	if err != nil {
		if r.Method != http.MethodPost {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	defer file.Close()

	source, binaryFileName, err := prepareAndSaveSource(file, language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// ソースファイルとバイナリファイルを削除
	defer os.Remove(source)
	defer os.Remove(binaryFileName)

	// コンパイル
	msg, err := runCompiler(language)
	if err != nil {
		log.Println(err.Error())
		msg := fmt.Sprintf("%v msg: %s", err, string(msg))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// google cloud storageにバイナリとソースファイルをアップロード
	rdm := generateRandomString(10)
	newfilename, err := uploadArtifacts(bucketName, source, binaryFileName, rdm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// バイナリのパスをレスポンスとして返す
	encodingFilename := url.QueryEscape(newfilename)
	w.Header().Set("Content-Disposition", "attachment; filename="+encodingFilename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	//_, err = w.Write(binary)
	//if err != nil {
	//http.Error(w, "Failed to write binary to response", http.StatusInternalServerError)
	//return
	//}
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
		log.Println("Failed to start server")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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

func uploadFileToGoogleCloudStorage(bucketName string, file, rdm string) (string, error) {
	if bucketName == "" {
		return "", fmt.Errorf("bucket name is not specified")
	}
	if file == "" {
		return "", fmt.Errorf("file name is not specified")
	}
	if rdm == "" {
		return "", fmt.Errorf("random string is not specified")
	}
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
