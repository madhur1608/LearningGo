package middleware

import (
	"net/http"
	"strings"

	"github.com/madhur1608/aditis-kitchen/customer-service/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "❌ Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			http.Error(w, "❌ Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		_, err := utils.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "❌ Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
