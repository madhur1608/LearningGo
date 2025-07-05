package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var adminEmail = "aditiskitchen17@gmail.com"
var jwtKey = []byte("aditi-secret")

func RequireAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "❌ Unauthorized - Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "❌ Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		email := (*claims)["email"]
		if email != adminEmail {
			http.Error(w, "❌ Unauthorized - Admin access required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
