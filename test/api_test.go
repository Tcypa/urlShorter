package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"url_shorter/inter/handler"
	"url_shorter/inter/storage"
	_ "url_shorter/inter/storage/memory"
)

var str = `memory`

func setupRouter() *mux.Router {
	storage.InitStorage(&str)

	router := mux.NewRouter()
	router.HandleFunc("/shorten", handler.ShortenURL).Methods("POST")
	router.HandleFunc("/{shortUrl}", handler.RedirectToOriginal).Methods("GET")
	return router
}

func TestShortenURL(t *testing.T) {
	router := setupRouter()

	payload := map[string]string{"url": "http://example.com"}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Проверка OK: получен статус %d", rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	if response["original_url"] != "http://example.com" {
		t.Errorf("получено %s", response["original_url"])
	}
	if response["short_url"] == "" {
		t.Error("short_url пустой")
	}
}

func TestInvalidRequest(t *testing.T) {
	router := setupRouter()

	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer([]byte(`{"invalid": "data"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Проверка BadRequest: получен статус %d", rr.Code)
	}
}

func TestRedirect(t *testing.T) {
	storage.InitStorage(&str)
	storePtr := storage.GetStorage()
	originalURL := "http://example.com"

	shortURL, err := (*storePtr).UrlInsert(originalURL)
	if err != nil {
		t.Fatalf("Ошибка при вставке URL: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/{shortUrl}", handler.RedirectToOriginal).Methods("GET")

	req, err := http.NewRequest("GET", "/"+shortURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Errorf("Получен статус %d", rr.Code)
	}
	if location := rr.Header().Get("Location"); location != originalURL {
		t.Errorf("Получен редирект %s, ожидался на %s", originalURL, location)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	storage.InitStorage(&str)
	router := mux.NewRouter()
	router.HandleFunc("/{shortUrl}", handler.RedirectToOriginal).Methods("GET")

	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Проверка NotFound: получен статус %d", rr.Code)
	}
}
