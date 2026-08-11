package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTest() {
	shortener = NewURLShortener()
}

func TestHandleShorten(t *testing.T) {
	setupTest()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "успешное сокращение URL",
			requestBody: ShortenRequest{
				URL: "http://example.com",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp ShortenResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("не удалось распарсить ответ: %v", err)
				}
				if resp.OriginalURL != "http://example.com" {
					t.Errorf("ожидался original_url 'http://example.com', получен '%s'", resp.OriginalURL)
				}
				if len(resp.ShortURL) < 6 {
					t.Errorf("короткий URL: %s", resp.ShortURL)
				}
			},
		},
		{
			name:           "некорректный JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp ErrorResponse
				if err := json.Unmarshal(body, &errResp); err != nil {
					t.Fatalf("не удалось распарсить ответ: %v", err)
				}
				if errResp.Error == "" {
					t.Error("ожидалось сообщение об ошибке")
				}
			},
		},
		{
			name: "пустой URL",
			requestBody: ShortenRequest{
				URL: "",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp ErrorResponse
				if err := json.Unmarshal(body, &errResp); err != nil {
					t.Fatalf("не удалось распарсить ответ: %v", err)
				}
				if errResp.Error == "" {
					t.Error("ожидалось сообщение об ошибке для пустого URL")
				}
			},
		},
		{
			name: "невалидный URL",
			requestBody: ShortenRequest{
				URL: "not-a-url",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp ErrorResponse
				if err := json.Unmarshal(body, &errResp); err != nil {
					t.Fatalf("не удалось распарсить ответ: %v", err)
				}
				if errResp.Error == "" {
					t.Error("ожидалось сообщение об ошибке для невалидного URL")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("не удалось маршалировать тело запроса: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handleShorten(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.expectedStatus, w.Code)
			}
			tt.checkResponse(t, w.Body.Bytes())
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	setupTest()

	shortID, err := shortener.Shorten("http://example.com/redirect")
	if err != nil {
		t.Fatalf("не удалось сохранить тестовый URL: %v", err)
	}

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedURL    string
	}{
		{
			name:           "успешный редирект",
			url:            "/" + shortID,
			expectedStatus: http.StatusFound,
			expectedURL:    "http://example.com/redirect",
		},
		{
			name:           "несуществующий короткий URL",
			url:            "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedURL:    "",
		},
		{
			name:           "пустой путь",
			url:            "/",
			expectedStatus: http.StatusBadRequest,
			expectedURL:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			handleRedirect(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusFound && tt.expectedURL != "" {
				location := w.Header().Get("Location")
				if location != tt.expectedURL {
					t.Errorf("ожидался Location '%s', получен '%s'", tt.expectedURL, location)
				}
			}
		})
	}
}

func TestMethodValidation(t *testing.T) {
	setupTest()

	t.Run("POST метод для /shorten", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
		w := httptest.NewRecorder()
		handleShorten(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("ожидался статус %d для неверного метода, получен %d",
				http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("не-GET метод для редиректа", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test123", nil)
		w := httptest.NewRecorder()
		handleRedirect(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("ожидался статус %d для неверного метода, получен %d",
				http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestIntegration(t *testing.T) {
	setupTest()

	t.Run("полный цикл: сокращение и редирект", func(t *testing.T) {
		body := bytes.NewBufferString(`{"url":"http://example.com/test"}`)
		req := httptest.NewRequest(http.MethodPost, "/shorten", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handleShorten(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("не удалось сократить URL, статус: %d", w.Code)
		}

		var resp ShortenResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("не удалось распарсить ответ: %v", err)
		}

		req = httptest.NewRequest(http.MethodGet, "/"+resp.ShortURL, nil)
		w = httptest.NewRecorder()
		handleRedirect(w, req)

		if w.Code != http.StatusFound {
			t.Errorf("ожидался статус редиректа %d, получен %d", http.StatusFound, w.Code)
		}

		location := w.Header().Get("Location")
		if location != "http://example.com/test" {
			t.Errorf("ожидался Location 'http://example.com/test', получен '%s'", location)
		}
	})
}
