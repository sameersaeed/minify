package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"minify/internal/limiter"
	"minify/internal/models"
	"minify/internal/services"
	"minify/internal/utils"

	"github.com/gorilla/mux"
)

type URLHandler struct {
	urlService       *services.URLService
	analyticsService *services.AnalyticsService
	limiter          *limiter.Limiter
}

func NewURLHandler(urlService *services.URLService, analyticsService *services.AnalyticsService, limiter *limiter.Limiter) *URLHandler {
	return &URLHandler{
		urlService:       urlService,       // handles db operations for URLs
		analyticsService: analyticsService, // records clicks and analytics
		limiter:          limiter,          // limits requests for each user
	}
}

// MinifyURL decodes the request, validates the given URL, generates a short code
// and returns the Minified URL
func (h *URLHandler) MinifyURL(w http.ResponseWriter, r *http.Request) {
	log.Println("[MinifyURL] Request received")
	var req models.MinifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[MinifyURL] Failed to decode request:", err)
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)

		return
	}
	log.Println("[MinifyURL] Original URL:", req.URL)

	ip := utils.GetClientIP(r)
	var (
		key string
		cfg limiter.RateConfig
	)

	// set rate limit based on whether requester is logged into an account
	if req.UserID != nil {
		key = fmt.Sprintf("user:%d", *req.UserID)
		cfg = limiter.Rates.Authenticated
	} else {
		key = fmt.Sprintf("ip:%s", ip)
		cfg = limiter.Rates.Anonymous
	}

	// enforce rate limit
	if !h.limiter.Allow(key, cfg) {
		log.Printf("[MinifyURL] Rate limit exceed for key %s", key)
		utils.JSONError(w, "Rate limit exceeded - please try again later.", http.StatusTooManyRequests)

		return
	}

	// validate url
	if !utils.IsValidURL(req.URL) {
		log.Println("[MinifyURL] Invalid URL format:", req.URL)
		utils.JSONError(w, "Invalid URL format", http.StatusBadRequest)

		return
	}

	// shorten (minify) url
	url, err := h.urlService.MinifyURL(req.URL, req.UserID)
	if err != nil {
		log.Println("[MinifyURL] Service failed:", err)
		utils.JSONError(w, "Failed to minify URL", http.StatusInternalServerError)

		return
	}
	log.Println("[MinifyURL] Short URL created:", url.ShortCode)

	baseURL := utils.GetBaseURL(r)
	response := models.MinifyResponse{
		ShortURL:    baseURL + "/" + url.ShortCode,
		OriginalURL: url.OriginalURL,
		ShortCode:   url.ShortCode,
	}

	utils.JSONResponse(w, response, http.StatusCreated)
	log.Println("[MinifyURL] Response sent")
}

// RedirectURL looks up the original URL by it's short code, increments click count (for metrics),
// records analytics, and redirects to the original URL
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]
	log.Println("[RedirectURL] Short code:", shortCode)

	if shortCode == "" {
		log.Println("[RedirectURL] Empty short code")
		http.NotFound(w, r)

		return
	}

	url, err := h.urlService.GetURLByShortCode(shortCode)
	if err != nil {
		log.Println("[RedirectURL] URL not found:", err)
		http.NotFound(w, r)

		return
	}

	go func() {
		if err := h.urlService.IncrementClickCount(url.ID); err != nil {
			log.Println("[RedirectURL] Failed to increment click count:", err)
		}
		h.analyticsService.RecordClick(url.ID, r.UserAgent(), utils.GetClientIP(r))
	}()

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

// GetUserURLs fetches all URLs associated with the given user ID
func (h *URLHandler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	log.Println("[GetUserURLs] User ID query param:", userIDStr)

	if userIDStr == "" {
		utils.JSONError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.JSONError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	urls, err := h.urlService.GetUserURLs(userID)
	if err != nil {
		log.Println("[GetUserURLs] Failed to get URLs:", err)
		utils.JSONError(w, "Failed to get URLs", http.StatusInternalServerError)

		return
	}

	utils.JSONResponse(w, urls, http.StatusOK)
	log.Println("[GetUserURLs] URLs returned:", len(urls))
}
