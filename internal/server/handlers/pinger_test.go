package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestPingHandler(t *testing.T) {

	type want struct {
		code int
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "PingHandler_No_DB_Connection",
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	router := chi.NewRouter()
	router.Get("/ping", PingHandler(Conf, TestLogger))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.code, result.StatusCode)
		})
	}
}

// --- Бенчмарки ---

func BenchmarkPingHandler(b *testing.B) {
	// Настраиваем MemMode для бенчмарка (предполагаем ошибку, т.к. нет БД)
	originalMode := Conf.MemMode
	Conf.MemMode = "memstore" // PingDB вернет ошибку
	defer func() { Conf.MemMode = originalMode }()

	router := chi.NewRouter()
	router.Get("/ping", PingHandler(Conf, TestLogger))

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}
