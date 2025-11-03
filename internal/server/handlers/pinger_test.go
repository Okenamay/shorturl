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
