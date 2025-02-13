package test

import (
	"math/rand"
	"testing"
	"url_shorter/internal/shorter"
)

func generateRandomURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := rand.Intn(10) + 5
	url := "https://random.site/"
	for i := 0; i < length; i++ {
		url += string(charset[rand.Intn(len(charset))])
	}
	return url
}

func TestShorter(t *testing.T) {
	mockURL := []string{
		"http://example.com",
		"https://google.com",
		"https://github.com",
		"https://golang.org",
	}
	for i := 0; i < 50; i++ {
		url := generateRandomURL()
		mockURL = append(mockURL, url)
	}
	generated := make(map[string]string)

	for _, url := range mockURL {
		shortURL := shorter.UrlShorter(url)

		if len(shortURL) < 10 {
			t.Errorf("%s малая длина:%d", shortURL, len(shortURL))
		}

		if secondShort := shorter.UrlShorter(url); shortURL != secondShort {
			t.Logf("Разные результаты для одного URL: %s -> %s и %s", url, shortURL, secondShort)
		}

		if existing, exists := generated[shortURL]; exists {
			t.Errorf("Коллизия: %s и %s сократились одинаково: %s", existing, url, shortURL)
		}

		generated[shortURL] = url
	}
}
