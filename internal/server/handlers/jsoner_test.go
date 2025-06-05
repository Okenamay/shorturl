package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONHandler(t *testing.T) {
	conf := config.ParseFlags()

	memstorage.URLStore = make(map[string]string)
	originalURL := "https://topdeck.ru/"
	result, shortID := urlmaker.ProcessURL(conf, originalURL)
	memstorage.URLStore[shortID] = originalURL

	type want struct {
		code        int
		response    JSONResponse
		contentType string
	}

	type request struct {
		method string
		url    string
		body   JSONRequest
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "JSONHandler_Correct_Method",
			request: request{
				method: http.MethodPost,
				url:    "/api/shorten",
				body:   JSONRequest{URL: originalURL},
			},
			want: want{
				code:        201,
				response:    JSONResponse{Result: result},
				contentType: "application/json",
			},
		},
		{
			name: "JSONHandler_Incorrect_Method",
			request: request{
				method: http.MethodGet,
				url:    "/api/shorten",
				body:   JSONRequest{},
			},
			want: want{
				code:        405,
				contentType: "",
			},
		},
	}

	router := chi.NewRouter()
	router.Post("/api/shorten", JSONHandler(conf))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(router)
			defer ts.Close()

			body, _ := json.Marshal(tt.request.body)
			request := httptest.NewRequest(tt.request.method, tt.request.url, nil)
			request.Body = io.NopCloser(bytes.NewReader(body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.code != http.StatusBadRequest && tt.want.code != http.StatusMethodNotAllowed {
				newURL, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				assert.NotEmpty(t, newURL)
			}
		})
	}
}
