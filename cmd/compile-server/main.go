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

	// ソースコードの読み込み
	src, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 一時ディレクトリ作成
	tmpDir, err := os.MkdirTemp("", "go-compiler-")
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create temp dir", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpDir)

	// ソースコードを一時ファイルに保存
	srcFile := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(srcFile, src, 0644)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to write source to file", http.StatusInternalServerError)
		return
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
