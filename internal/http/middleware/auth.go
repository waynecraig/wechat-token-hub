package middleware

import (
	"net/http"
	"strings"

	"github.com/waynecraig/wechat-token-hub/internal/auth"
)

// a middleware function that checks the authorization header for a valid JWT token
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// check that the authorization header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		// get the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// parse and validate the token
		if err := auth.VerifyJwtToken(tokenString); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// call the next handler
		next.ServeHTTP(w, r)
	})
}
