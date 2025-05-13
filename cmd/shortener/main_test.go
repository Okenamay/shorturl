package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	type request struct {
		method string
		url    string
		body   []byte
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "ShortenHandler Correct Method",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("https://scryfall.com"),
			},
			want: want{
				code:        201,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "ShortenHandler Incorrect Method",
			request: request{
				method: http.MethodGet,
				url:    "/",
				body:   []byte("https://www.mtggoldfish.com/"),
			},
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name: "ShortenHandler Incorrect Scheme",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("ftp://tcgplayer.com/"),
			},
			want: want{
				code:        422,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "ShortenHandler Incorrect URL",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("hilmar.v.petursson@ccpgames.com"),
			},
			want: want{
				code:        422,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.url, bytes.NewReader(tt.request.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortenHandler)
			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.code != http.StatusBadRequest {
				newURL, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				assert.NotEmpty(t, newURL)
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	config.ParseFlags()
	URLStore = make(map[string]string)
	originalURL := "https://topdeck.ru/"
	hash := md5.New()
	io.WriteString(hash, originalURL)
	shortID := hex.EncodeToString(hash.Sum(nil))[:config.Cfg.ShortIDLen]
	URLStore[shortID] = originalURL

	type want struct {
		code        int
		response    string
		contentType string
	}

	type request struct {
		method string
		url    string
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "RedirectHandler Correct Method",
			request: request{
				method: http.MethodGet,
				url:    "/" + shortID,
			},
			want: want{
				code:        307,
				response:    originalURL,
				contentType: "text/plain",
			},
		},
		{
			name: "RedirectHandler Wrong Method",
			request: request{
				method: http.MethodPost,
				url:    "/" + shortID,
			},
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/{id}", RedirectHandler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(router)
			defer ts.Close()

			fullURL := ts.URL + tt.request.url

			parsedURL, err := url.Parse(fullURL)
			require.NoError(t, err)

			request := httptest.NewRequest(tt.request.method, fullURL, nil)
			request.URL = parsedURL

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.code != http.StatusBadRequest {
				require.Equal(t, originalURL, result.Header.Get("Location"))
			}
		})
	}
}
