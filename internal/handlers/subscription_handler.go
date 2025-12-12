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

type SubscriptionHandler struct {
	service *services.SubscriptionService
}

func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// Subscribe subscribes a user to one or more packages
func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
		return
	}

	var subscribeReq models.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&subscribeReq); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(subscribeReq); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	response, err := h.service.Subscribe(userID, subscribeReq.PackageIDs)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, utils.ResponseTypeResource, response, response.Message)
}

// GetActive retrieves all active subscriptions for the authenticated user
func (h *SubscriptionHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
		return
	}

	subscriptions, err := h.service.GetActiveSubscriptions(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to retrieve active subscriptions")
		return
	}

	response := map[string]interface{}{
		"subscriptions": subscriptions,
		"total":         len(subscriptions),
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeCollection, response, "Active subscriptions retrieved successfully")
}

// GetHistory retrieves all subscriptions (active and expired) for the authenticated user
func (h *SubscriptionHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
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

	subscriptions, err := h.service.GetAllSubscriptions(userID, limit, offset)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to retrieve subscription history")
		return
	}

	// Get total count
	count, err := h.service.CountByUserID(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to count subscriptions")
		return
	}

	response := map[string]interface{}{
		"subscriptions": subscriptions,
		"total":         count,
		"limit":         limit,
		"offset":        offset,
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeCollection, response, "Subscription history retrieved successfully")
}

// CheckAccess checks if user has access to specific asset class and duration type
func (h *SubscriptionHandler) CheckAccess(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
		return
	}

	var checkReq models.CheckAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&checkReq); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(checkReq); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	response, err := h.service.CheckAccess(userID, checkReq.AssetClass, checkReq.DurationType)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to check access")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeResource, response, "Access check completed successfully")
}

