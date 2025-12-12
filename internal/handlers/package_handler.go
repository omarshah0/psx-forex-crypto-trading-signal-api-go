package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/services"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type PackageHandler struct {
	service *services.PackageService
}

func NewPackageHandler(service *services.PackageService) *PackageHandler {
	return &PackageHandler{service: service}
}

// GetAll retrieves all active packages (public)
func (h *PackageHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 100
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

	packages, err := h.service.GetAll(true, limit, offset) // activeOnly = true
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to retrieve packages")
		return
	}

	// Get total count
	count, err := h.service.Count(true)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to count packages")
		return
	}

	response := map[string]interface{}{
		"packages": packages,
		"total":    count,
		"limit":    limit,
		"offset":   offset,
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeCollection, response, "Packages retrieved successfully")
}

// GetByID retrieves a single package by ID
func (h *PackageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid package ID")
		return
	}

	pkg, err := h.service.GetByID(id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, utils.ErrorTypeNotFound, "Package not found")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeResource, pkg, "Package retrieved successfully")
}

// Create creates a new package (admin only)
func (h *PackageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var packageCreate models.PackageCreate
	if err := json.NewDecoder(r.Body).Decode(&packageCreate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(packageCreate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	pkg, err := h.service.Create(&packageCreate)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to create package")
		return
	}

	utils.SendSuccess(w, http.StatusCreated, utils.ResponseTypeResource, pkg, "Package created successfully")
}

// Update updates a package (admin only)
func (h *PackageHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid package ID")
		return
	}

	var packageUpdate models.PackageUpdate
	if err := json.NewDecoder(r.Body).Decode(&packageUpdate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(packageUpdate); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	pkg, err := h.service.Update(id, &packageUpdate)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to update package")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeResource, pkg, "Package updated successfully")
}

// Delete deletes a package (admin only)
func (h *PackageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid package ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to delete package")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "Package deleted successfully")
}

