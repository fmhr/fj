package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("starting server...")
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("handling request...")
	w.Write([]byte("hello"))
}
