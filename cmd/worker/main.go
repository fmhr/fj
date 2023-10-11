package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Configの受け取り
	configPart := r.FormValue("config")
	var config fj.Config
	err = json.Unmarshal([]byte(configPart), &config)
	if err != nil {
		http.Error(w, "Failed to unmarshal config", http.StatusBadRequest)
		return
	}

	// バイナリの受け取り
	file, _, err := r.FormFile("binary")
	if err != nil {
		http.Error(w, "Failed to get the binary", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "uploaded-binary-*")
	if err != nil {
		http.Error(w, "Failed to create a temp file", http.StatusInternalServerError)
		return
	}
	tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		errmsg := fmt.Sprint("Failed to copy the binary to the temp file", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	// language
	language := r.FormValue("language")
	if language == "" {
		http.Error(w, "Language not specified", http.StatusBadRequest)
		return
	}
	// seed
	seedString := r.FormValue("seed")
	seedInt, err := strconv.Atoi(seedString)
	if err != nil {
		log.Printf("seed not specified %s", seedString)
		http.Error(w, "Seed not specified", http.StatusBadRequest)
		return
	}
	out, err := exexute(&config, seedInt)
	if err != nil {
		http.Error(w, "Failed to execute", http.StatusInternalServerError)
		return
	}
	// json
	jsonData, err := json.Marshal(out)
	if err != nil {
		http.Error(w, "Failed to marshal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/upload", handler)
	http.ListenAndServe(":8080", nil)
}
