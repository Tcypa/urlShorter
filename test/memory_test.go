package test

import (
	"testing"

	"url_shorter/inter/storage/memory"
)

func TestMemStorage(t *testing.T) {
	store := memory.Init()
	defer store.Close()

	originalURL := "http://example.com"
	shortURL, err := store.UrlInsert(originalURL)
	if err != nil {
		t.Fatalf("UrlInsert вернул ошибку: %v", err)
	}
	if shortURL == "" {
		t.Error("UrlInsert вернул пустой Url")
	}

	shortURL2, err := store.UrlInsert(originalURL)
	if err != nil {
		t.Fatalf("Вторая вставка UrlInsert вернул ошибку: %v", err)
	}
	if shortURL != shortURL2 {
		t.Error("Повторная вставка вернула другое значение")
	}

	readURL, err := store.UrlRead(shortURL)
	if err != nil {
		t.Fatalf("UrlRead вернул ошибку: %v", err)
	}
	if readURL != originalURL {
		t.Errorf("Оригинальный Url: %s не равен прочитанному, получено %s", originalURL, readURL)
	}
}
