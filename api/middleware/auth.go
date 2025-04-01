package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ashil-poojary/banking-ledger-service/utils"
	"github.com/go-redis/redis/v8"
)

// AuthMiddleware checks if a user is authenticated
func AuthMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				utils.SendResponse(w, http.StatusUnauthorized, false, "", nil, "Invalid Authorization")
				return
			}

			// Extract token from "Bearer <token>"
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate JWT and extract user_id
			userID, err := utils.ParseJWT(tokenString)
			if err != nil {
				utils.SendResponse(w, http.StatusUnauthorized, false, "", nil, "Invalid Authorization")
				return
			}

			// Check if session exists in Redis using user_id
			exists, err := redisClient.Exists(ctx, userID).Result()
			if err != nil {
				utils.SendResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to check session")
				return
			}
			if exists == 0 {
				utils.SendResponse(w, http.StatusUnauthorized, false, "", nil, "Invalid Authorization")
				return
			}

			// Store user_id in context for further request handling
			ctx = context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
