package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecretKey = "super-secret-key-for-testing"

// TestJWTAuthRoundTrip проверяет, что мы можем успешно создать токен и сразу
// же его проверить
func TestJWTAuthRoundTrip(t *testing.T) {
	conf := &config.Cfg{
		TokenExpiry:      24,
		AuthorizationKey: testSecretKey,
	}
	testUserID := "test-user-123"

	tokenString, err := buildJWTString(conf, testUserID)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	parsedUserID, err := GetUserID(conf, tokenString)
	require.NoError(t, err)
	assert.Equal(t, testUserID, parsedUserID)

	t.Run("Invalid Signature", func(t *testing.T) {
		wrongConf := &config.Cfg{AuthorizationKey: "wrong-key"}
		_, err := GetUserID(wrongConf, tokenString)
		assert.Error(t, err)
	})

	t.Run("Malformed Token", func(t *testing.T) {
		_, err := GetUserID(conf, "not.a.real.token")
		assert.Error(t, err)
	})

	t.Run("Expired Token", func(t *testing.T) {
		expiredConf := &config.Cfg{
			TokenExpiry:      -1, // -1 час
			AuthorizationKey: testSecretKey,
		}

		expiredToken, err := buildJWTString(expiredConf, testUserID)
		require.NoError(t, err)

		_, err = GetUserID(expiredConf, expiredToken)
		assert.Error(t, err)
	})
}

// TestAuthenticatorMiddleware проверяет сам HTTP middleware
func TestAuthenticatorMiddleware(t *testing.T) {
	conf := &config.Cfg{
		TokenExpiry:      24,
		AuthorizationKey: testSecretKey,
	}

	// dummyHandler - это внутренний хендлер, который будет вызван, если
	// middleware пропустит запрос. Он проверяет контекст и пишет userID в тело
	// ответа
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := CheckAuth(r)
		if !ok {
			http.Error(w, "UserID not found in context", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userID))
	})

	authMiddleware := Authenticator(conf)(dummyHandler)

	t.Run("No Cookie - New User Created", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		authMiddleware.ServeHTTP(rr, req)

		result := rr.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, result.StatusCode)

		cookie := result.Cookies()[0]
		assert.Equal(t, "token", cookie.Name)
		assert.NotEmpty(t, cookie.Value)

		body := rr.Body.String()
		assert.Equal(t, 36, len(body), "Expected body to be a UUID string of length 36")
	})

	t.Run("Valid Cookie - Existing User", func(t *testing.T) {
		testUserID := "existing-user-id"
		validToken, err := buildJWTString(conf, testUserID)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: validToken,
		})

		rr := httptest.NewRecorder()
		authMiddleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Equal(t, testUserID, rr.Body.String())
	})

	t.Run("Invalid Cookie - Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "this-is-not-a-valid-jwt",
		})

		rr := httptest.NewRecorder()
		authMiddleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Unauthorized\n", rr.Body.String())
	})
}

func TestCheckAuth(t *testing.T) {
	t.Run("Context With UserID", func(t *testing.T) {
		expectedUserID := "user-from-context"
		ctx := context.WithValue(context.Background(), UserIDContextKey, expectedUserID)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		userID, ok := CheckAuth(req)
		assert.True(t, ok)
		assert.Equal(t, expectedUserID, userID)
	})

	t.Run("Context Without UserID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil) // Пустой контекст
		userID, ok := CheckAuth(req)
		assert.False(t, ok)
		assert.Empty(t, userID)
	})

	t.Run("Context With Wrong Type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDContextKey, 12345)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		userID, ok := CheckAuth(req)
		assert.False(t, ok)
		assert.Empty(t, userID)
	})
}

// --- Бенчмарки ---

var (
	benchConf = &config.Cfg{
		TokenExpiry:      24,
		AuthorizationKey: testSecretKey,
	}
	benchUserID   = "benchmark-user-id"
	benchResult   string
	benchCtx      context.Context
	benchToken, _ = buildJWTString(benchConf, benchUserID)
)

func BenchmarkBuildJWTString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = buildJWTString(benchConf, benchUserID)
	}
}

func BenchmarkGetUserID(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = GetUserID(benchConf, benchToken)
	}
}

func BenchmarkAuthenticatorMiddleware(b *testing.B) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		benchCtx = r.Context() // Сохраняем контекст, чтобы избежать оптимизации
		w.WriteHeader(http.StatusOK)
	})
	authMiddleware := Authenticator(benchConf)(dummyHandler)

	validCookieReq := httptest.NewRequest(http.MethodGet, "/", nil)
	validCookieReq.AddCookie(&http.Cookie{Name: "token", Value: benchToken})

	b.Run("NewUser (No Cookie)", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Необходимо создавать новый запрос в цикле, т.к. ServeHTTP может
			// его модифицировать
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			authMiddleware.ServeHTTP(rr, req)
		}
	})

	b.Run("ValidUser (Valid Cookie)", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			rr := httptest.NewRecorder()
			// req можно переиспользовать, т.к. он не меняется
			authMiddleware.ServeHTTP(rr, validCookieReq)
		}
	})
}
