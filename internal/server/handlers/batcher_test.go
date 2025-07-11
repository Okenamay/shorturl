package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestBatchHandlerTransaction(t *testing.T) {
	memstorage.Store = memstorage.NewURLMap()

	requestPayload := []RequestEntry{
		{CorrelationID: "1", OriginalURL: "https://google.com"},
		{CorrelationID: "2", OriginalURL: "https://yandex.ru"},
		{CorrelationID: "3", OriginalURL: "https://bing.com"},
	}
	body, _ := json.Marshal(requestPayload)

	type want struct {
		code        int
		contentType string
		respLen     int
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "BatchHandler_Successful_Transaction",
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
				respLen:     3,
			},
		},
	}

	router := chi.NewRouter()
	router.Post("/api/shorten/batch", BatchHandlerTransaction(Conf))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			var responsePayload []ResponseEntry
			err := json.NewDecoder(result.Body).Decode(&responsePayload)
			require.NoError(t, err)
			require.Len(t, responsePayload, tt.want.respLen)
			require.Equal(t, "1", responsePayload[0].CorrelationID)
			require.NotEmpty(t, responsePayload[0].ShortURL)
		})
	}
}
