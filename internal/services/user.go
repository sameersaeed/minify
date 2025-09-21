package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"minify/internal/models"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser inserts a new user into the db, with their provided password being hashed
func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	log.Printf("[UserService] Creating user: %s\n", username)

	// hash before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[UserService] Failed to hash password:", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	var user models.User
	err = s.db.QueryRow(query, username, email, string(hashedPassword)).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		log.Println("[UserService] DB insert failed:", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.Username = username
	user.Email = email
	log.Println("[UserService] User created with ID:", user.ID)
	return &user, nil
}

// AuthenticateUser fetches a user by their username for login
func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	log.Println("[UserService] Authenticating user:", username)
	query := `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[UserService] User", user.Username, "not found")
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
