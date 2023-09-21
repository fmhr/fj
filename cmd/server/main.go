package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fmhr/fj"
)

func main() {
	log.Println("starting server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on port " + port)
	http.HandleFunc("/", handler)
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
	config.GenPath = "gen"
	config.VisPath = "vis"
	config.TesterPath = "tester"
	reactiveString := r.URL.Query().Get("reactive")
	if reactiveString == "" {
		config.Reactive = false
	} else {
		config.Reactive = true
	}
	// GEN
	fj.Gen(&config, seed)
	// RUN
	rtn, err := run(&config, seed)
	if err != nil {
		log.Println("run error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		return fj.ReactiveRun(cfg, seed)
	}
	return fj.RunVis(cfg, seed)
}
