package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sergey/work-track-backend/internal/middleware"
	"github.com/sergey/work-track-backend/internal/models"
	"github.com/sergey/work-track-backend/internal/repository"
	"github.com/sergey/work-track-backend/internal/service"
)

// TrackItemHandler handles track item endpoints
type TrackItemHandler struct {
	trackItemService *service.TrackItemService
}

// NewTrackItemHandler creates a new track item handler
func NewTrackItemHandler(trackItemService *service.TrackItemService) *TrackItemHandler {
	return &TrackItemHandler{
		trackItemService: trackItemService,
	}
}

// ListTrackItems retrieves all track items for the authenticated user
func (h *TrackItemHandler) ListTrackItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check for date range query parameters
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var items []models.TrackItem
	var err error

	if startDate != "" && endDate != "" {
		// Get items by date range
		items, err = h.trackItemService.GetTrackItemsByDateRange(r.Context(), userID, startDate, endDate)
	} else {
		// Get all items
		items, err = h.trackItemService.GetUserTrackItems(r.Context(), userID)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

// GetTrackItemsSummary retrieves aggregated shift totals for the authenticated user
func (h *TrackItemHandler) GetTrackItemsSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	if startDate == "" || endDate == "" {
		respondWithError(w, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	summary, err := h.trackItemService.GetTrackItemsSummaryByDateRange(r.Context(), userID, startDate, endDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, summary)
}

// CreateTrackItem creates a new track item
func (h *TrackItemHandler) CreateTrackItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateTrackItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.trackItemService.CreateTrackItem(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, item)
}

// GetTrackItem retrieves a specific track item
func (h *TrackItemHandler) GetTrackItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	itemID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid track item ID")
		return
	}

	item, err := h.trackItemService.GetTrackItem(r.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrTrackItemNotFound) {
			respondWithError(w, http.StatusNotFound, "Track item not found")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			respondWithError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

// UpdateTrackItem updates a track item
func (h *TrackItemHandler) UpdateTrackItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	itemID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid track item ID")
		return
	}

	var req models.UpdateTrackItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.trackItemService.UpdateTrackItem(r.Context(), userID, itemID, &req)
	if err != nil {
		if errors.Is(err, repository.ErrTrackItemNotFound) {
			respondWithError(w, http.StatusNotFound, "Track item not found")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			respondWithError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

// DeleteTrackItem deletes a track item
func (h *TrackItemHandler) DeleteTrackItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	itemID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid track item ID")
		return
	}

	err = h.trackItemService.DeleteTrackItem(r.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrTrackItemNotFound) {
			respondWithError(w, http.StatusNotFound, "Track item not found")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			respondWithError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
