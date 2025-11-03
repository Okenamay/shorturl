package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	type request struct {
		method string
		url    string
		body   JSONRequest
	}

	tests := []struct {
		name    string
		setup   func()
		request request
		want    want
	}{
		{
			name: "JSONHandler_Correct_Method_New_URL",
			setup: func() {
				memstorage.Store = memstorage.NewURLMap()
			},
			request: request{
				method: http.MethodPost,
				url:    "/api/shorten",
				body:   JSONRequest{URL: "https://new-url.com"},
			},
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name: "JSONHandler_Conflict_When_URL_Exists",
			setup: func() {
				memstorage.Store = memstorage.NewURLMap()
				existingURL := "https://existing-url.com"
				_, shortID := urlmaker.ProcessURL(Conf, existingURL)
				memstorage.Store.Set(shortID, existingURL)
			},
			request: request{
				method: http.MethodPost,
				url:    "/api/shorten",
				body:   JSONRequest{URL: "https://existing-url.com"},
			},
			want: want{
				code:        http.StatusConflict,
				contentType: "application/json",
			},
		},
		{
			name: "JSONHandler_Incorrect_Method",
			setup: func() {
				memstorage.Store = memstorage.NewURLMap()
			},
			request: request{
				method: http.MethodGet,
				url:    "/api/shorten",
				body:   JSONRequest{},
			},
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
	}

	router := chi.NewRouter()
	router.Post("/api/shorten", JSONHandler(Conf, TestLogger, nil))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			ts := httptest.NewServer(router)
			defer ts.Close()

			body, _ := json.Marshal(tt.request.body)
			request := httptest.NewRequest(tt.request.method, tt.request.url, bytes.NewReader(body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.code == http.StatusCreated || tt.want.code == http.StatusConflict {
				var resp JSONResponse
				err := json.NewDecoder(result.Body).Decode(&resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Result)
			}
		})
	}
}

// --- Бенчмарки ---

func BenchmarkJSONHandler(b *testing.B) {
	router := chi.NewRouter()
	router.Post("/api/shorten", JSONHandler(Conf, TestLogger, nil))

	newURLBody, _ := json.Marshal(JSONRequest{URL: "https://new-url-for-bench.com"})
	newURLBytes := newURLBody

	// Готовим тело для конфликтующего URL
	conflictURL := "https://existing-url-for-bench.com"
	conflictBody, _ := json.Marshal(JSONRequest{URL: conflictURL})
	conflictURLBytes := conflictBody

	b.Run("New URL", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			memstorage.Store = memstorage.NewURLMap()
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(newURLBytes))
			rr := httptest.NewRecorder()
			b.StartTimer()

			router.ServeHTTP(rr, req)
		}
	})

	b.Run("Conflict URL", func(b *testing.B) {
		memstorage.Store = memstorage.NewURLMap()
		_, shortID := urlmaker.ProcessURL(Conf, conflictURL)
		memstorage.Store.Set(shortID, conflictURL)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Не нужно сбрасывать хранилище, т.к. мы тестируем конфликт
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(conflictURLBytes))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
		}
	})
}
