package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// ExtractUserID extracts the user ID from a Bearer token in the Authorization header
func ExtractUserID(r *http.Request) (string, error) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	log.Println("[DEBUG] Extracting token from Authorization header:", authHeader)

	if authHeader == "" {
		log.Println("[ERROR] Authorization token missing")
		return "", errors.New("authorization token required")
	}

	// Extract token from "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == "" {
		return "", errors.New("bearer token missing")
	}

	// Get secret key from environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret" // Fallback for testing
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("[ERROR] Unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		log.Println("[ERROR] Invalid token:", err)
		return "", errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("[ERROR] Invalid token claims")
		return "", errors.New("invalid token claims")
	}

	// Extract user ID
	userIDStr, ok := claims["user_id"].(string)

	if !ok {
		log.Println("[ERROR] user_id not found in token")
		return "", errors.New("user_id not found in token")
	}

	return userIDStr, nil
}

// parseJWT extracts the UserID from the token
func ParseJWT(tokenString string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}

// GenerateJWT creates a signed token// GenerateJWT creates a signed token with UserID
func GenerateJWT(userID string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret"
	}

	claims := jwt.MapClaims{
		"user_id": userID, // Store UserID instead of username
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
