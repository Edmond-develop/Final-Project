package api

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil || cookie.Value == "" {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			jwtParse, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(pass), nil
			})
			if err != nil || !jwtParse.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			claims, ok := jwtParse.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			if claims["password"] != HashPassword(pass) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
