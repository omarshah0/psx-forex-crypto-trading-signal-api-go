package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/omarshah0/rest-api-with-social-auth/internal/middleware"
	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/services"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type TradingSignalHandler struct {
	service *services.TradingSignalService
}

func NewTradingSignalHandler(service *services.TradingSignalService) *TradingSignalHandler {
	return &TradingSignalHandler{service: service}
}

// GetAll retrieves trading signals visible to the authenticated user
func (h *TradingSignalHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	// Get signals visible to this user based on their subscriptions
	signals, err := h.service.GetSignalsForUser(userID, limit, offset)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to retrieve trading signals")
		return
	}

	// Get total count for this user
	count, err := h.service.CountForUser(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to count trading signals")
		return
	}

	response := map[string]interface{}{
		"signals": signals,
		"total":   count,
		"limit":   limit,
		"offset":  offset,
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeCollection, response, "Trading signals retrieved successfully")
}

// GetAllAdmin retrieves all trading signals (admin only)
func (h *TradingSignalHandler) GetAllAdmin(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	signals, err := h.service.GetAll(limit, offset)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to retrieve trading signals")
		return
	}

	// Get total count
	count, err := h.service.Count()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to count trading signals")
		return
	}

	response := map[string]interface{}{
		"signals": signals,
		"total":   count,
		"limit":   limit,
		"offset":  offset,
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeCollection, response, "Trading signals retrieved successfully")
}

// GetByID retrieves a single trading signal by ID
func (h *TradingSignalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid signal ID")
		return
	}

	signal, err := h.service.GetByID(id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, utils.ErrorTypeNotFound, "Trading signal not found")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeResource, signal, "Trading signal retrieved successfully")
}

// Create creates a new trading signal (admin only)
func (h *TradingSignalHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
		return
	}

	var signalCreate models.TradingSignalCreate
	if err := json.NewDecoder(r.Body).Decode(&signalCreate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(signalCreate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	signal, err := h.service.Create(&signalCreate, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to create trading signal")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, utils.ResponseTypeResource, signal, "Trading signal created successfully")
}

// Update updates a trading signal (admin only)
func (h *TradingSignalHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid signal ID")
		return
	}

	var signalUpdate models.TradingSignalUpdate
	if err := json.NewDecoder(r.Body).Decode(&signalUpdate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(signalUpdate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	signal, err := h.service.Update(id, &signalUpdate)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to update trading signal")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeResource, signal, "Trading signal updated successfully")
}

// Delete deletes a trading signal (admin only)
func (h *TradingSignalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid signal ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to delete trading signal")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "Trading signal deleted successfully")
}
