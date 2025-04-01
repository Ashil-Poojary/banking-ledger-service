package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/ashil-poojary/banking-ledger-service/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountHandler handles account-related requests
type AccountHandler struct {
	DB *gorm.DB
}

// NewAccountHandler initializes a new AccountHandler
func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{DB: db}
}

// CreateAccount handles account creation
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserID(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Unauthorized", nil, err.Error())
		return
	}

	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, "Invalid request payload", nil, err.Error())
		return
	}

	account.UserID = userID
	account.ID = uuid.New()

	if err := h.DB.Create(&account).Error; err != nil {
		log.Println("Failed to create account:", err)
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to create account", nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusCreated, true, "Account created successfully", account, "")
}

// GetAccount retrieves the authenticated user's account
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserID(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Unauthorized", nil, err.Error())
		return
	}

	var account models.Account
	if err := h.DB.Where("user_id = ?", userID).First(&account).Error; err != nil {
		utils.SendResponse(w, http.StatusNotFound, false, "Account not found", nil, "")
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Account retrieved successfully", account, "")
}

func (h *AccountHandler) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserID(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Unauthorized", nil, err.Error())
		return
	}

	var accounts []models.Account
	if err := h.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to retrieve accounts", nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Accounts retrieved successfully", accounts, "")
}

// UpdateAccount updates the authenticated user's account
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserID(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Unauthorized", nil, err.Error())
		return
	}

	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, "Invalid request payload", nil, err.Error())
		return
	}

	if err := h.DB.Model(&models.Account{}).Where("user_id = ?", userID).Updates(account).Error; err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to update account", nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Account updated successfully", account, "")
}

// DeleteAccount deletes the authenticated user's account
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserID(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Unauthorized", nil, err.Error())
		return
	}

	if err := h.DB.Delete(&models.Account{}, "user_id = ?", userID).Error; err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to delete account", nil, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Account deleted successfully", nil, "")
}
