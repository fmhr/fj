package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	// バイナリの受け取り
	file, _, err := r.FormFile("binary")
	if err != nil {
		errmsg := fmt.Sprint("Failed to get the binary:", err.Error())
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}
	defer file.Close()

	//	tmpFile, err := os.CreateTemp("", "uploaded-binary-*")
	//if err != nil {
	//errmsg := fmt.Sprint("Failed to create a temp file", err.Error())
	//http.Error(w, errmsg, http.StatusInternalServerError)
	//return
	//}
	//defer tmpFile.Close()

	binaryPath, err := os.Create(config.BinaryName)
	if err != nil {
		errmsg := fmt.Sprint("Failed to create a binary file", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(binaryPath, file)
	if err != nil {
		errmsg := fmt.Sprint("Failed to copy the binary to the temp file", err.Error())
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

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
	http.HandleFunc("/upload", handler)
	http.ListenAndServe(":8080", nil)
}
