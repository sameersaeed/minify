package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"minify/internal/config"
	"minify/internal/models"
	"minify/internal/services"
	"minify/internal/utils"
)

type UserHandler struct {
	userService *services.UserService	// handles db operations on users
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser handles user registration
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("[CreateUser] Request received")
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[CreateUser] Failed to decode request:", err)
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("[CreateUser] Username: %s, Email: %s\n", req.Username, req.Email)

	if err := utils.ValidateStruct(req); err != nil {
		log.Println("[CreateUser] Validation failed:", err)
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Println("[CreateUser] Service error:", err)
		utils.JSONError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Println("[CreateUser] User created successfully:", user.ID)
	utils.JSONResponse(w, user, http.StatusCreated)
}

// LoginUser handles user login and issuing JWT auth tokens
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("[LoginUser] Request received")
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[LoginUser] Failed to decode request:", err)
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		log.Println("[LoginUser] Authentication failed:", err)
		utils.JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	log.Println("[LoginUser] Authentication successful for user:", user.ID)

	token, err := h.generateJWTToken(user.ID, user.Username)
	if err != nil {
		log.Println("[LoginUser] Failed to generate JWT:", err)
		utils.JSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}
	utils.JSONResponse(w, response, http.StatusOK)
	log.Println("[LoginUser] Response sent with token")
}

// generateJWTToken creates a signed JWT token for the provided user ID and username
func (h *UserHandler) generateJWTToken(userID int, username string) (string, error) {
	cfg := config.Load()
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),	// 1 week
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}
