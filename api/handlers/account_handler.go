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

	// Extract account number from query parameter
	accountNumber := r.URL.Query().Get("account_number")
	if accountNumber == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, "Account number is required", nil, "")
		return
	}

	// Ensure the account belongs to the user
	var account models.Account
	if err := h.DB.Where("account_number = ? AND user_id = ?", accountNumber, userID).First(&account).Error; err != nil {
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

// UpdateAccount updates a specific account belonging to the authenticated user
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserID(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Unauthorized", nil, err.Error())
		return
	}

	accountNumber := r.URL.Query().Get("account_number")
	if accountNumber == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, "Account number is required", nil, "")
		return
	}

	var updateData models.Account
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, "Invalid request payload", nil, err.Error())
		return
	}

	// Ensure the account belongs to the user and update it
	result := h.DB.Model(&models.Account{}).
		Where("account_number = ? AND user_id = ?", accountNumber, userID).
		Updates(updateData)

	if result.Error != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to update account", nil, result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, false, "Account not found or no changes applied", nil, "")
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Account updated successfully", updateData, "")
}

// DeleteAccount deletes a specific account belonging to the authenticated user
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserID(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, "Unauthorized", nil, err.Error())
		return
	}

	accountNumber := r.URL.Query().Get("account_number")
	if accountNumber == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, "Account number is required", nil, "")
		return
	}

	// Ensure the account belongs to the user before deleting
	result := h.DB.Where("account_number = ? AND user_id = ?", accountNumber, userID).Delete(&models.Account{})

	if result.Error != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "Failed to delete account", nil, result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, false, "Account not found", nil, "")
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Account deleted successfully", nil, "")
}
