package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

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

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	language := r.FormValue("language")
	if language == "" {
		http.Error(w, "Language not specified", http.StatusBadRequest)
		return
	}

	// 一時ディレクトリ作成
	tmpDir, err := os.MkdirTemp("", "go-compiler-")
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create temp dir", http.StatusInternalServerError)
		return
	}
	// defer os.RemoveAll(tmpDir)
	srcFile := filepath.Join(tmpDir, "main.go")
	srcBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read the uploaded file", http.StatusInternalServerError)
		return
	}
	err = os.WriteFile(srcFile, srcBytes, 0644)
	if err != nil {
		http.Error(w, "Failed to write the uploaded file to disk", http.StatusInternalServerError)
		return
	}
	if language == "go" {
		log.Println("go file")
	}

	// Goでコンパイル
	outFile := filepath.Join(tmpDir, "a.out")
	cmd := exec.Command("go", "build", "-o", outFile, srcFile)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "Compiration failed", http.StatusBadRequest)
		return
	}

	// バイナリの読み込み
	binary, err := os.ReadFile(outFile)
	if err != nil {
		http.Error(w, "Failed to read compiled binary", http.StatusInternalServerError)
		return
	}

	// バイナリをレスポンスとして返す
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(binary)
}

func main() {
	http.HandleFunc("/compile", compileHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":"+port, nil)
}
