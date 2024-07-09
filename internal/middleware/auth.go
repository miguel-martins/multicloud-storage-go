package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/miguel-martins/multicloud-storage-go/internal/util"
)

type contextKey string

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]
		jwt, err := util.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// TODO Get user associated with username and send to contenxt
		const usernameKey contextKey = "username"
		username := jwt.Username
		// TODO Future referernce perform additional checks here based on the claims, like checking user roles, etc.
		r = r.WithContext(context.WithValue(r.Context(), usernameKey, username))
		next.ServeHTTP(w, r)
	})
}
