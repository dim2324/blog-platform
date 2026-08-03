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

// Shorten создает короткий идентификатор для URL
func (us *URLShortener) Shorten(originalURL string) (string, error) {
	// Валидация URL
	if !isValidURL(originalURL) {
		return "", fmt.Errorf("невалидный URL: %s", originalURL)
	}

	// Генерация уникального короткого ID
	shortID := generateShortID()

	// Проверка уникальности и повторная генерация при необходимости
	us.mu.Lock()
	defer us.mu.Unlock()

	// Убедимся, что ID уникален
	for {
		if _, exists := us.urls[shortID]; !exists {
			break
		}
		shortID = generateShortID()
	}

	us.urls[shortID] = originalURL
	return shortID, nil
}

// GetOriginal возвращает оригинальный URL по короткому ID
func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	originalURL, exists := us.urls[shortID]
	if !exists {
		return "", fmt.Errorf("короткий URL не найден: %s", shortID)
	}

	return originalURL, nil
}

// generateShortID генерирует случайный короткий идентификатор длиной 6-8 символов
func generateShortID() string {
	b := make([]byte, 6) // 6 байт = 8 символов в base64
	rand.Read(b)
	encoded := base64.URLEncoding.EncodeToString(b)
	// Обрезаем до 8 символов и убираем символы, которые могут быть неудобны в URL
	encoded = strings.TrimRight(encoded, "=")
	if len(encoded) > 8 {
		encoded = encoded[:8]
	}
	return encoded
}

// isValidURL проверяет корректность URL
func isValidURL(str string) bool {
	if str == "" {
		return false
	}

	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	// Проверяем, что схема - http или https
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Проверяем, что хост не пустой
	if u.Host == "" {
		return false
	}

	return true
}
