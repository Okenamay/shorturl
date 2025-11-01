package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/audit"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	testConf    *config.Cfg
	testLogger  *zap.SugaredLogger
	testAuditor *audit.Auditor
	testRouter  http.Handler
)

// TestMain настраивает общие зависимости для всех тестов в этом пакете.
func TestMain(m *testing.M) {
	testConf = config.InitConfig()
	if testConf.AuthorizationKey == "" {
		testConf.AuthorizationKey = "test-secret-key"
	}

	testLogger = zap.NewNop().Sugar()

	testAuditor = audit.NewAuditor("", "", testLogger)

	testRouter = NewRouter(testConf, testLogger, testAuditor)

	os.Exit(m.Run())
}

// TestRoutes проверяет, что роутер правильно маршрутизирует запросы
func TestRoutes(t *testing.T) {
	ts := httptest.NewServer(testRouter)
	defer ts.Close()

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "GET nonexistent route",
			method:     http.MethodGet,
			path:       "/nonexistent/route",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Method not allowed for /",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Method not allowed for /api/shorten",
			method:     http.MethodGet,
			path:       "/api/shorten",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Auth middleware check (unauthorized)",
			method:     http.MethodGet,
			path:       "/api/user/urls",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+tt.path, nil)
			assert.NoError(t, err)

			res, err := client.Do(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

// --- Бенчмарки ---

func BenchmarkRoutes(b *testing.B) {
	// Тестовые сценарии, аналогичные TestRoutes
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET /ping",
			method: http.MethodGet,
			path:   "/ping",
		},
		{
			name:   "GET nonexistent route",
			method: http.MethodGet,
			path:   "/nonexistent/route",
		},
		{
			name:   "Method not allowed for /",
			method: http.MethodGet,
			path:   "/",
		},
		{
			name:   "Auth middleware check (unauthorized)",
			method: http.MethodGet,
			path:   "/api/user/urls",
		},
	}

	// Запускаем бенчмарк для каждого сценария
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Создаем запрос и рекордер на каждой итерации
				req := httptest.NewRequest(tt.method, tt.path, nil)
				rr := httptest.NewRecorder()

				// Выполняем запрос напрямую через testRouter.ServeHTTP, что
				// быстрее, чем запускать httptest.Server
				testRouter.ServeHTTP(rr, req)
			}
		})
	}
}
