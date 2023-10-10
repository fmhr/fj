package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
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
		http.Error(w, "Failed to save the binary", http.StatusInternalServerError)
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
	out, err := exexute(tmpFile.Name(), language, seedInt)
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
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}
