package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hawkerd/privateinstruction/internal/auth"
)

func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract the token
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// parse the token
		claims, err := auth.ParseJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// extract the user ID
		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// add user ID to the request context
		ctx := context.WithValue(r.Context(), "userID", uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
