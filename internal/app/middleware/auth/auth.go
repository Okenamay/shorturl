package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

func buildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUserID(tokenString string) (string, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(SecretKey), nil
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

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			newUserID := uuid.New().String()
			tokenString, buildErr := buildJWTString(newUserID)
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

		userID, err := getUserID(cookie.Value)
		if err != nil {

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
