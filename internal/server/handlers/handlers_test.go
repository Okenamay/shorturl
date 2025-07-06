package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper"
	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Conf *config.Cfg

func TestMain(m *testing.M) {
	if err := logger.InitLogger(); err != nil {
		logger.Zap.Fatalw(err.Error(), "Main", "Start logger")
	}

	Conf = config.InitConfig()

	os.Exit(m.Run())
}

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
			name: "ShortenHandler_Correct_Method",
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
			name: "ShortenHandler_Incorrect_Method",
			request: request{
				method: http.MethodGet,
				url:    "/",
				body:   []byte("https://www.mtggoldfish.com/"),
			},
			want: want{
				code:        405,
				response:    "",
				contentType: "",
			},
		},
		{
			name: "ShortenHandler_Incorrect_Scheme",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("ftp://tcgplayer.com/"),
			},
			want: want{
				code:        400,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "ShortenHandler_Incorrect_URL",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("hilmar.v.petursson@ccpgames.com"),
			},
			want: want{
				code:        400,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "ShortenHandler_Conflict_URL_Exists",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("https://google.com/maps"),
			},
			want: want{
				code:        http.StatusConflict,
				contentType: "text/plain",
			},
		},
	}

	router := chi.NewRouter()
	router.With(gzipper.Decompressor, gzipper.Compressor).Post("/", ShortenHandler(Conf))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memstorage.Store = memstorage.NewURLMap()

			if tt.name == "ShortenHandler_Conflict_When_URL_Exists" {
				_, shortID := urlmaker.ProcessURL(Conf, string(tt.request.body))
				_, err := memselect.StorePair(Conf, shortID, string(tt.request.body))
				require.NoError(t, err)
			}

			request := httptest.NewRequest(tt.request.method, tt.request.url, bytes.NewReader(tt.request.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.code == http.StatusCreated || tt.want.code == http.StatusConflict {
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				assert.NotEmpty(t, string(body))
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	memstorage.Store = memstorage.NewURLMap()
	originalURL := "https://topdeck.ru/"
	_, shortID := urlmaker.ProcessURL(Conf, originalURL)
	memstorage.Store.Set(shortID, originalURL)

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
			name: "RedirectHandler_Correct_Method",
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
			name: "RedirectHandler_Wrong_Method",
			request: request{
				method: http.MethodPost,
				url:    "/" + shortID,
			},
			want: want{
				code:        405,
				response:    "",
				contentType: "",
			},
		},
	}

	router := chi.NewRouter()
	router.With(gzipper.Decompressor, gzipper.Compressor).Get("/{id}", RedirectHandler(Conf))

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
			defer result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if result.StatusCode != http.StatusMethodNotAllowed {
				require.Equal(t, originalURL, result.Header.Get("Location"))
			}
		})
	}
}
