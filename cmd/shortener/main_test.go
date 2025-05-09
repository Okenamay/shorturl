package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
	url := "https://topdeck.ru/"

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
				url:    "/",
			},
			want: want{
				code:        307,
				response:    url,
				contentType: "text/plain",
			},
		},
		{
			name: "RedirectHandler Wrong Method",
			request: request{
				method: http.MethodPost,
				url:    "/",
			},
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(url)))
	w := httptest.NewRecorder()
	h := http.HandlerFunc(ShortenHandler)
	h(w, req)

	result := w.Result()
	newURL, _ := io.ReadAll(result.Body)
	result.Body.Close()

	fmt.Println(string(newURL))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, string(newURL), nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(RedirectHandler)
			h(w, request)

			result := w.Result()
			result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.code != http.StatusBadRequest {
				require.Equal(t, url, result.Header.Get("Location"))
			}
		})
	}
}
