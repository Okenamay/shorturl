package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Okenamay/shorturl.git/internal/app/middleware/auth"
	"github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper"
	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/Okenamay/shorturl.git/internal/worker"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Conf *config.Cfg
var TestLogger *zap.SugaredLogger

func TestMain(m *testing.M) {
	TestLogger, err := logger.InitLogger()
	if err != nil {
		TestLogger.Fatalw("Tests stopped - start logger FAIL", "error", err)
	}
	defer TestLogger.Sync()

	Conf = config.InitConfig()
	Conf.MemMode = "memstore"

	worker.DeleteChan = make(chan worker.DeleteTask, 128)

	os.Exit(m.Run())
}

func TestBatchDeleter(t *testing.T) {
	router := chi.NewRouter()
	router.Delete("/api/user/urls", BatchDeleter(Conf))

	t.Run("BatchDeleter_Accepts_Request", func(t *testing.T) {
		shortIDs := []string{"a", "b", "c"}
		body, _ := json.Marshal(shortIDs)

		request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))

		ctx := context.WithValue(request.Context(), auth.UserIDContextKey, "test-user-for-delete")
		request = request.WithContext(ctx)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		result := w.Result()
		defer result.Body.Close()

		require.Equal(t, http.StatusAccepted, result.StatusCode)

		select {
		case task := <-worker.DeleteChan:
			assert.Equal(t, "test-user-for-delete", task.UserID)
			assert.Equal(t, shortIDs, task.ShortIDs)
		case <-time.After(100 * time.Millisecond):
			t.Error("handler did not send a task to the delete channel in time")
		}
	})
}

func TestRedirectHandler(t *testing.T) {
	// Пустая затычка теста на доступ к удалённому URL:
	t.Run("RedirectHandler_Requesting_Deleted_URL", func(t *testing.T) {
		originalMode := Conf.MemMode
		Conf.MemMode = "postgres"
		defer func() { Conf.MemMode = originalMode }()

		w := httptest.NewRecorder()

		w.WriteHeader(http.StatusGone)
		result := w.Result()
		defer result.Body.Close()
		require.Equal(t, http.StatusGone, result.StatusCode)
	})

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
	router.With(
		gzipper.Decompressor(TestLogger),
		gzipper.Compressor(TestLogger),
	).Get("/{id}", RedirectHandler(Conf, TestLogger))

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

func TestShortenHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	type request struct {
		method string
		url    string
		body   []byte
	}

	tests := []struct {
		name    string
		setup   func()
		request request
		want    want
	}{
		{
			name: "ShortenHandler_Correct_Method_New_URL",
			setup: func() {
				memstorage.Store = memstorage.NewURLMap()
			},
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("https://new-url.com"),
			},
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name: "ShortenHandler_Conflict_When_URL_Exists",
			setup: func() {
				memstorage.Store = memstorage.NewURLMap()
				existingURL := "https://existing-url.com"
				_, shortID := urlmaker.ProcessURL(Conf, existingURL)
				memstorage.Store.Set(shortID, existingURL)
			},
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("https://existing-url.com"),
			},
			want: want{
				code:        http.StatusConflict,
				contentType: "text/plain",
			},
		},
		{
			name: "ShortenHandler_Incorrect_Method",
			setup: func() {
				memstorage.Store = memstorage.NewURLMap()
			},
			request: request{
				method: http.MethodGet,
				url:    "/",
				body:   []byte("https://www.mtggoldfish.com/"),
			},
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name: "ShortenHandler_Incorrect_Scheme",
			setup: func() {
				memstorage.Store = memstorage.NewURLMap()
			},
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("ftp://tcgplayer.com/"),
			},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "ShortenHandler_Incorrect_URL",
			setup: func() {
				memstorage.Store = memstorage.NewURLMap()
			},
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   []byte("hilmar.v.petursson@ccpgames.com"),
			},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	router := chi.NewRouter()
	router.Post("/", ShortenHandler(Conf, TestLogger))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
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

func TestUserURLsHandler(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/api/user/urls", UserURLsHandler(Conf, TestLogger))

	t.Run("UserURLsHandler_Unauthorized", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)

		result := w.Result()
		defer result.Body.Close()

		require.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})

	t.Run("UserURLsHandler_No_Content", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		ctx := context.WithValue(request.Context(), auth.UserIDContextKey, "test-user-id")
		request = request.WithContext(ctx)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)

		result := w.Result()
		defer result.Body.Close()

		require.Equal(t, http.StatusNoContent, result.StatusCode)
	})
}
