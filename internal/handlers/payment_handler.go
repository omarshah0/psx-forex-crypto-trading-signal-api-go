package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/omarshah0/rest-api-with-social-auth/internal/middleware"
	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/services"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// GetHistory retrieves payment history for the authenticated user
func (h *PaymentHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
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

	payments, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to retrieve payment history")
		return
	}

	// Get total count
	count, err := h.service.CountByUserID(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to count payments")
		return
	}

	response := map[string]interface{}{
		"payments": payments,
		"total":    count,
		"limit":    limit,
		"offset":   offset,
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeCollection, response, "Payment history retrieved successfully")
}

// RecordPayment manually records a payment (admin only, dummy for now)
func (h *PaymentHandler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	var paymentCreate models.PaymentCreate
	if err := json.NewDecoder(r.Body).Decode(&paymentCreate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(paymentCreate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	payment, err := h.service.Create(&paymentCreate)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to record payment")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, utils.ResponseTypeResource, payment, "Payment recorded successfully")
}

