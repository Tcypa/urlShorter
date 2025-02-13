package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"url_shorter/inter/storage"
)

type URLMapping struct {
	Original string `json:"original_url"`
	Short    string `json:"short_url"`
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" || !IsUrl(req.URL) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	store := *storage.GetStorage()

	shortUrl, err := store.UrlInsert(req.URL)
	if err != nil {
		log.Printf("Can't create short url: %v", err)
	}

	response := URLMapping{
		Original: req.URL,
		Short:    "http://localhost:8080/" + shortUrl,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["shortUrl"]
	store := *storage.GetStorage()
	origUrl, err := store.UrlRead(shortUrl)

	if err != nil || origUrl == "" {
		log.Printf("Short URL not found: %s", shortUrl)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	log.Printf("Redirecting: %s -> %s", shortUrl, origUrl)
	http.Redirect(w, r, origUrl, http.StatusFound)
}
