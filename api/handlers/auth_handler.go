package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/ashil-poojary/banking-ledger-service/utils"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles user authentication
type AuthHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(db *gorm.DB, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		DB:    db,
		Redis: redisClient,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, "Invalid request", nil, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to hash password")
		return
	}
	user.Password = string(hashedPassword)

	// Save user in DB
	if err := h.DB.Create(&user).Error; err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "User registration failed", nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusCreated, true, "User registered successfully", nil, "")
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var dbUser models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, "Invalid request", nil, err.Error())
		return
	}

	// Find user in DB
	if err := h.DB.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Invalid credentials", nil, err.Error())
		return
	}
	log.Printf("User registered")
	log.Print(dbUser)

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Invalid credentials", nil, err.Error())
		return
	}

	// Generate JWT token with UserID
	token, err := utils.GenerateJWT(dbUser.ID.String())
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to generate token", nil, err.Error())
		return
	}

	// Store session in Redis with UserID as key
	ctx := context.Background()
	err = h.Redis.Set(ctx, dbUser.ID.String(), token, 24*time.Hour).Err()
	if err != nil {
		log.Println("Failed to store session in Redis:", err)
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to store session", nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Login successful", map[string]string{"token": token}, "")
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	token := r.Header.Get("Authorization")

	if token == "" {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Authorization token required", nil, "No token provided")
		return
	}

	// Extract UserID from token
	userID, err := utils.ParseJWT(token)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Invalid token", nil, err.Error())
		return
	}

	// Remove session from Redis using UserID
	err = h.Redis.Del(ctx, userID).Err()
	if err != nil {
		log.Println("Failed to delete session from Redis:", err)
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to logout", nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Logged out successfully", nil, "")
}
