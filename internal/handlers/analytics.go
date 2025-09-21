package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"minify/internal/services"
	"minify/internal/utils"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService	// handles fetching analytics data
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetOverview fetches the overall analytics and returns them as JSON
func (h *AnalyticsHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	log.Println("[Analytics] GetOverview request received")
	stats, err := h.analyticsService.GetOverview()
	if err != nil {
		log.Println("[Analytics] Failed to get overview:", err)
		utils.JSONError(w, "Failed to get overview stats", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, stats, http.StatusOK)
	log.Println("[Analytics] Overview stats sent")
}

// GetPopularURLs fetches the top N URLs (default 10) by clicks and returns them as JSON
func (h *AnalyticsHandler) GetPopularURLs(w http.ResponseWriter, r *http.Request) {
	log.Println("[Analytics] GetPopularURLs request received")
	limitStr := r.URL.Query().Get("limit")
	limit := 10 
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	urls, err := h.analyticsService.GetPopularURLs(limit)
	if err != nil {
		log.Println("[Analytics] Failed to get popular URLs:", err)
		utils.JSONError(w, "Failed to get popular URLs", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, urls, http.StatusOK)
	log.Printf("[Analytics] Sent top %d popular URLs\n", limit)
}

// GetTimeframeStats fetches analytics for a specific period (hour, day, month, year)
func (h *AnalyticsHandler) GetTimeframeStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	period := vars["period"]
	log.Println("[Analytics] GetTimeframeStats request received for period:", period)

	validPeriods := map[string]bool{"hour": true, "day": true, "month": true, "year": true}
	if !validPeriods[period] {
		log.Println("[Analytics] Invalid period:", period)
		utils.JSONError(w, "Invalid period. Use: hour, day, month, year", http.StatusBadRequest)
		return
	}

	stats, err := h.analyticsService.GetTimeframeStats(period)
	if err != nil {
		log.Println("[Analytics] Failed to get timeframe stats:", err)
		utils.JSONError(w, "Failed to get timeframe stats", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, stats, http.StatusOK)
	log.Println("[Analytics] Timeframe stats sent for period:", period)
}
