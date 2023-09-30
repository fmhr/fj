package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fmhr/fj"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	log.Println("starting server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on port " + port)
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("handling request...")
	// seed
	seedString := r.URL.Query().Get("seed")
	if seedString == "" {
		http.Error(w, "no seed specified", http.StatusBadRequest)
		return
	}
	// seed
	seed, err := strconv.Atoi(seedString)
	if err != nil {
		http.Error(w, "invalid seed specified", http.StatusBadRequest)
		return
	}
	log.Println("seed:", seed)

	var config fj.Config
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read body error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = toml.Unmarshal(body, &config)
	if err != nil {
		log.Println("toml unmarshal error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	config.Cmd = os.Getenv("EXECUTE_COMMAND")
	config.GenPath = "tools/target/release/gen"
	config.VisPath = "tools/target/release/vis"
	config.TesterPath = "tools/target/release/tester"
	reactiveString := r.URL.Query().Get("reactive")
	if reactiveString == "" || reactiveString == "false" {
		config.Reactive = false
	} else {
		config.Reactive = true
	}
	// GEN
	err = fj.Gen(&config, seed)
	if err != nil {
		log.Println("gen error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// RUN
	rtn, err := run(&config, seed)
	if err != nil {
		log.Println("run error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("rtn:", rtn)
	jsonData, err := json.Marshal(rtn)
	if err != nil {
		log.Println("json marshal error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func run(cfg *fj.Config, seed int) (map[string]float64, error) {
	if cfg.Reactive {
		log.Println("reactive mode")
		return fj.ReactiveRun(cfg, seed)
	}
	log.Println("normal mode")
	return fj.RunVis(cfg, seed)
}
