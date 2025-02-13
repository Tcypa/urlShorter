package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	cfg "url_shorter/config"
	api "url_shorter/internal/handler"
	stg "url_shorter/internal/storage"
)

func main() {
	storageType := flag.String("stgType", "memory", "Choose storage type (postgres/memory)")
	flag.Parse()
	cfg.LoadConfig("config/config.yaml")
	stg.InitStorage(storageType)

	r := mux.NewRouter()
	r.HandleFunc("/shorten", api.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortUrl}", api.RedirectToOriginal).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
