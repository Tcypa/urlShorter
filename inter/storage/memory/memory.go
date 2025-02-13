package memory

import (
	"errors"
	"fmt"
	"sync"
	"url_shorter/inter/shorter"
)

type MemStorage struct {
	mu              sync.RWMutex
	origToShortUrls map[string]string
	shortToOrigUrls map[string]string
}

func Init() *MemStorage {
	fmt.Println("Memory storage init")
	return &MemStorage{origToShortUrls: make(map[string]string), shortToOrigUrls: make(map[string]string)}
}

func (m *MemStorage) Close() {
	//не требует закрытия
}

func (m *MemStorage) UrlInsert(origUrl string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	shortUrl, exists := m.origToShortUrls[origUrl]
	if exists {
		return shortUrl, nil
	} else {
		shortUrl = shorter.UrlShorter(origUrl)
		m.origToShortUrls[origUrl] = shortUrl
		m.shortToOrigUrls[shortUrl] = origUrl

		return shortUrl, nil
	}
}

func (m *MemStorage) UrlRead(shortUrl string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	origUrl, exists := m.shortToOrigUrls[shortUrl]
	if !exists {
		return "", errors.New("short URL not found")
	}

	return origUrl, nil
}
