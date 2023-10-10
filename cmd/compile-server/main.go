package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
	compileCmd := r.FormValue("compileCmd")
	if compileCmd == "" {
		http.Error(w, "Language not specified", http.StatusBadRequest)
		return
	}
	// ソースファイル名
	srcFileName := r.FormValue("srcFile")
	if srcFileName == "" {
		http.Error(w, "Source file not specified", http.StatusBadRequest)
		return
	}
	// バイナリファイル名
	binaryFileName := r.FormValue("binaryFile")

	srcBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read the uploaded file", http.StatusInternalServerError)
		return
	}
	err = os.WriteFile(srcFileName, srcBytes, 0644)
	if err != nil {
		http.Error(w, "Failed to write the uploaded file to disk", http.StatusInternalServerError)
		return
	}

	cmds := strings.Fields(compileCmd)

	cmd := exec.Command(cmds[0], cmds[1:]...)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "Compiration failed", http.StatusBadRequest)
		return
	}

	// バイナリの読み込み
	binary, err := os.ReadFile(binaryFileName)
	if err != nil {
		http.Error(w, "Failed to read compiled binary", http.StatusInternalServerError)
		return
	}
	// ソースファイルとバイナリファイルを削除
	defer os.Remove(srcFileName)
	defer os.Remove(binaryFileName)

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
