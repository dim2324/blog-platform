package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type URLShortener struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

func (us *URLShortener) Shorten(originalURL string) (string, error) {
	if !isValidURL(originalURL) {
		return "", fmt.Errorf("невалидный URL: %s", originalURL)
	}

	shortID := generateShortID()

	us.mu.Lock()
	defer us.mu.Unlock()

	for {
		if _, exists := us.urls[shortID]; !exists {
			break
		}
		shortID = generateShortID()
	}

	us.urls[shortID] = originalURL
	return shortID, nil
}

func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	originalURL, exists := us.urls[shortID]
	if !exists {
		return "", fmt.Errorf("короткий URL не найден: %s", shortID)
	}

	return originalURL, nil
}

func generateShortID() string {
	b := make([]byte, 6)
	rand.Read(b)
	encoded := base64.URLEncoding.EncodeToString(b)
	encoded = strings.TrimRight(encoded, "=")
	if len(encoded) > 8 {
		encoded = encoded[:8]
	}
	return encoded
}

func isValidURL(str string) bool {
	// Явная проверка на пробелы – тест ожидает false для URL с пробелами
	if strings.Contains(str, " ") {
		return false
	}

	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}
