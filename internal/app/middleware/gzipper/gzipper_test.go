package gzipper

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// mustDecompress считывает gzipped-байт и возвращает распакованную строку
func mustDecompress(t *testing.T, data []byte) string {
	t.Helper()

	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	require.NoError(t, err)
	defer gzReader.Close()

	decompressedBody, err := io.ReadAll(gzReader)
	require.NoError(t, err)

	return string(decompressedBody)
}

// mustCompress сжимает строку и возвращает gzipped-байты
func mustCompress(t *testing.T, data string) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	_, err := gzWriter.Write([]byte(data))
	require.NoError(t, err)
	err = gzWriter.Close()
	require.NoError(t, err)

	return &buf
}

func TestCompressorMiddleware(t *testing.T) {
	testLogger := zap.NewNop().Sugar()
	const responseBody = `{"hello": "world"}`

	// nextHandler - это хендлер, который вызывается после middleware,
	// устанавливает Content-Type и пишет тело ответа
	nextHandler := func(contentType string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(responseBody))
		})
	}

	tests := []struct {
		name                 string
		acceptEncodingHeader string
		contentType          string
		wantCompressed       bool
	}{
		{
			name:                 "Gzip accepted, valid Content-Type (json)",
			acceptEncodingHeader: "gzip",
			contentType:          "application/json",
			wantCompressed:       true,
		},
		{
			name:                 "Gzip accepted, valid Content-Type (html)",
			acceptEncodingHeader: "gzip",
			contentType:          "text/html; charset=utf-8",
			wantCompressed:       true,
		},
		{
			name:                 "Gzip accepted, invalid Content-Type (text/plain)",
			acceptEncodingHeader: "gzip",
			contentType:          "text/plain",
			wantCompressed:       false,
		},
		{
			name:                 "Gzip NOT accepted",
			acceptEncodingHeader: "identity", // без сжатия
			contentType:          "application/json",
			wantCompressed:       false,
		},
		{
			name:                 "No Accept-Encoding header",
			acceptEncodingHeader: "",
			contentType:          "application/json",
			wantCompressed:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.acceptEncodingHeader != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncodingHeader)
			}

			rr := httptest.NewRecorder()

			handler := Compressor(testLogger)(nextHandler(tt.contentType))
			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			if tt.wantCompressed {
				assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))
				decompressed := mustDecompress(t, rr.Body.Bytes())
				assert.Equal(t, responseBody, decompressed)
			} else {
				assert.Empty(t, rr.Header().Get("Content-Encoding"))
				assert.Equal(t, responseBody, rr.Body.String())
			}
		})
	}
}

func TestDecompressorMiddleware(t *testing.T) {
	testLogger := zap.NewNop().Sugar()
	const requestBody = `{"ping": "pong"}`

	// nextHandler - это эхо-сервер, который читает тело запроса и пишет его
	// в тело ответа. Если декомпрессия успешна, он вернет распакованное тело
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	handler := Decompressor(testLogger)(nextHandler)

	t.Run("No Content-Encoding (plain text)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, requestBody, rr.Body.String())
	})

	t.Run("With gzip Content-Encoding (compressed)", func(t *testing.T) {
		compressedBody := mustCompress(t, requestBody)

		req := httptest.NewRequest(http.MethodPost, "/", compressedBody)
		req.Header.Set("Content-Encoding", "gzip")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, requestBody, rr.Body.String())
	})

	t.Run("With gzip Content-Encoding but corrupt body", func(t *testing.T) {
		// Отправляем обычный текст, но с заголовком gzip
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("this is not gzipped"))
		req.Header.Set("Content-Encoding", "gzip")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		// gzip.NewReader вернет ошибку, и middleware должен вернуть 500
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

// --- Бенчмарки ---

var (
	benchLogger      = zap.NewNop().Sugar()
	benchCompressReq = func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		return req
	}()
	benchDecompressReqBody = mustCompress(&testing.T{}, `{"field_one": "value_one", "field_two": "value_two"}`)
	benchDecompressReq     = func() *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(benchDecompressReqBody.Bytes()))
		req.Header.Set("Content-Encoding", "gzip")
		return req
	}()
)

func BenchmarkCompressorMiddleware(b *testing.B) {
	handler := Compressor(benchLogger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"field_one": "value_one", "field_two": "value_two"}`))
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, benchCompressReq)
	}
}

func BenchmarkDecompressorMiddleware(b *testing.B) {
	handler := Decompressor(benchLogger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Читаем тело, чтобы r.Body был полностью обработан
		io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Необходимо создать новый ридер для тела, т.к. оно читается 1 раз
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(benchDecompressReqBody.Bytes()))
		req.Header.Set("Content-Encoding", "gzip")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
