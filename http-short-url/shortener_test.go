package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"валидный HTTP URL", "http://example.com", true},
		{"валидный HTTPS URL", "https://google.com/search?q=test", true},
		{"валидный URL с портом", "http://localhost:8080/path", true},
		{"невалидный URL без схемы", "not-a-url", false},
		{"пустая строка", "", false},
		{"URL без хоста", "http://", false},
		{"FTP URL", "ftp://example.com", false},
		{"валидный URL с поддоменом", "https://sub.domain.com/path", true},
		{"URL с пробелами", "http://example.com/path with spaces", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("isValidURL(%q) = %v, ожидалось %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestURLShortener_Shorten(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"валидный HTTP URL", "http://example.com", false},
		{"валидный HTTPS URL", "https://google.com/search?q=test", false},
		{"невалидный URL", "not-a-url", true},
		{"пустая строка", "", true},
		{"FTP URL", "ftp://example.com", true},
		{"валидный URL с длинным путем", "https://example.com/very/long/path/to/resource", false},
	}

	shortener := NewURLShortener()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := shortener.Shorten(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shorten() ошибка = %v, ожидалась ошибка = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(shortID) < 6 || len(shortID) > 8 {
					t.Errorf("короткий ID имеет неверную длину: %d, ожидалось 6-8", len(shortID))
				}
				originalURL, err := shortener.GetOriginal(shortID)
				if err != nil {
					t.Errorf("не удалось получить оригинальный URL: %v", err)
				}
				if originalURL != tt.url {
					t.Errorf("ожидался URL %q, получен %q", tt.url, originalURL)
				}
			}
		})
	}
}

func TestURLShortener_GetOriginal(t *testing.T) {
	shortener := NewURLShortener()

	testURL := "http://example.com"
	shortID, err := shortener.Shorten(testURL)
	if err != nil {
		t.Fatalf("не удалось сохранить тестовый URL: %v", err)
	}

	tests := []struct {
		name    string
		shortID string
		wantURL string
		wantErr bool
	}{
		{"существующий короткий ID", shortID, testURL, false},
		{"несуществующий короткий ID", "nonexistent", "", true},
		{"пустой короткий ID", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalURL, err := shortener.GetOriginal(tt.shortID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOriginal() ошибка = %v, ожидалась ошибка = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && originalURL != tt.wantURL {
				t.Errorf("GetOriginal() = %q, ожидалось %q", originalURL, tt.wantURL)
			}
		})
	}
}

func TestGenerateShortID(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := generateShortID()
		if len(id) < 6 || len(id) > 8 {
			t.Errorf("сгенерирован ID неверной длины: %d", len(id))
		}
		if ids[id] {
			t.Errorf("сгенерирован дублирующийся ID: %s", id)
		}
		ids[id] = true
	}
}

func TestURLShortener_Concurrent(t *testing.T) {
	shortener := NewURLShortener()
	const goroutines = 10
	errCh := make(chan error, goroutines)
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			url := fmt.Sprintf("http://example.com/path/%d", id)
			if _, err := shortener.Shorten(url); err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("ошибка при конкурентном сохранении: %v", err)
	}
}
