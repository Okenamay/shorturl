package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

func buildJWTString(conf *config.Cfg, userID string) (string, error) {
	tokenExp := time.Duration(conf.TokenExpiry) * time.Hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(conf.AuthorizationKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID экспортируем для использования в gRPC интерцепторе
func GetUserID(conf *config.Cfg, tokenString string) (string, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(conf.AuthorizationKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	return claims.UserID, nil
}

type contextKey string

// UserIDContextKey - это ключ, используемый для хранения ID пользователя в
// context.Context запроса.
const UserIDContextKey = contextKey("userID")

// Authenticator - это middleware, которое управляет аутентификацией
// пользователя. Если у пользователя нет "token" cookie:
// 1. Генерируется новый UUID пользователя.
// 2. Создается новый JWT-токен.
// 3. Токен устанавливается в "token" cookie.
// 4. UUID пользователя добавляется в контекст запроса.
//
// Если "token" cookie есть:
// 1. Токен валидируется.
// 2. В случае успеха, UserID из токена добавляется в контекст запроса.
// 3. В случае неудачи (неверная подпись, истекший срок), возвращается 401
// Unauthorized.
func Authenticator(conf *config.Cfg) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenExp := time.Duration(conf.TokenExpiry) * time.Hour

			cookie, err := r.Cookie("token")
			if err != nil {
				newUserID := uuid.New().String()
				tokenString, buildErr := buildJWTString(conf, newUserID)
				if buildErr != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:    "token",
					Value:   tokenString,
					Path:    "/",
					Expires: time.Now().Add(tokenExp),
				})

				ctx := context.WithValue(r.Context(), UserIDContextKey, newUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, err := GetUserID(conf, cookie.Value)
			if err != nil {

				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CheckAuth - вспомогательная функция для извлечения ID пользователя из
// context.Context запроса.
// Возвращает ID пользователя и true, если он был найден и аутентифицирован.
// Возвращает пустую строку и false, если пользователь не аутентифицирован.
func CheckAuth(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}
