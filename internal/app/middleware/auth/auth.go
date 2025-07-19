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

const TokenExp = time.Hour * 24
const SecretKey = "supersecretkey"

func buildJWTString(conf *config.Cfg, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(conf.AuthorizationKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUserID(conf *config.Cfg, tokenString string) (string, error) {
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

const UserIDContextKey = contextKey("userID")

func Authenticator(conf *config.Cfg) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
					Expires: time.Now().Add(TokenExp),
				})

				ctx := context.WithValue(r.Context(), UserIDContextKey, newUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, err := getUserID(conf, cookie.Value)
			if err != nil {

				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
