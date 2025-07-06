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

// func TestJSONHandler(t *testing.T) {
// 	memstorage.Store = memstorage.NewURLMap()
// 	originalURL := "https://topdeck.ru/"
// 	result, shortID := urlmaker.ProcessURL(Conf, originalURL)
// 	memstorage.Store.Set(shortID, originalURL)

// 	type want struct {
// 		code        int
// 		response    JSONResponse
// 		contentType string
// 	}

// 	type request struct {
// 		method string
// 		url    string
// 		body   JSONRequest
// 	}

// 	tests := []struct {
// 		name    string
// 		request request
// 		want    want
// 	}{
// 		{
// 			name: "JSONHandler_Correct_Method",
// 			request: request{
// 				method: http.MethodPost,
// 				url:    "/api/shorten",
// 				body:   JSONRequest{URL: originalURL},
// 			},
// 			want: want{
// 				code:        201,
// 				response:    JSONResponse{Result: result},
// 				contentType: "application/json",
// 			},
// 		},
// 		{
// 			name: "JSONHandler_Incorrect_Method",
// 			request: request{
// 				method: http.MethodGet,
// 				url:    "/api/shorten",
// 				body:   JSONRequest{},
// 			},
// 			want: want{
// 				code:        405,
// 				contentType: "",
// 			},
// 		},
// 	}

// 	router := chi.NewRouter()
// 	router.Post("/api/shorten", JSONHandler(Conf))

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ts := httptest.NewServer(router)
// 			defer ts.Close()

// 			body, _ := json.Marshal(tt.request.body)
// 			request := httptest.NewRequest(tt.request.method, tt.request.url, bytes.NewReader(body))
// 			w := httptest.NewRecorder()
// 			router.ServeHTTP(w, request)

// 			result := w.Result()
// 			defer result.Body.Close()

// 			require.Equal(t, tt.want.code, result.StatusCode)
// 			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

// 			if tt.want.code != http.StatusBadRequest && tt.want.code != http.StatusMethodNotAllowed {
// 				var resp JSONResponse
// 				err := json.NewDecoder(result.Body).Decode(&resp)
// 				require.NoError(t, err)
// 				assert.NotEmpty(t, resp.Result)
// 			}
// 		})
// 	}
// }

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
		setup   func() // Setup function to prepare the store for each specific test case
		request request
		want    want
	}{
		{
			name: "JSONHandler_Correct_Method_New_URL",
			setup: func() {
				// For a new URL, the store must be empty.
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
				// For a conflict test, pre-populate the store with the URL.
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
	router.Post("/api/shorten", JSONHandler(Conf))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the specific setup for this test case.
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
