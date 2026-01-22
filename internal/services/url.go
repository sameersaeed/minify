package services

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	"minify/internal/models"
)

type URLService struct {
	db *sql.DB
}

func NewURLService(db *sql.DB) *URLService {
	return &URLService{db: db}
}

// MinifyURL generates a unique short code and inserts the original URL into the db
func (s *URLService) MinifyURL(originalURL string, userID *int) (*models.URL, error) {
	shortCode, err := s.generateShortCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	for s.shortCodeExists(shortCode) {
		shortCode, err = s.generateShortCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate unique short code: %w", err)
		}
	}

	query := `
		INSERT INTO urls (short_code, original_url, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	var url models.URL
	err = s.db.QueryRow(query, shortCode, originalURL, userID).Scan(
		&url.ID,
		&url.CreatedAt,
		&url.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert URL: %w", err)
	}

	url.ShortCode = shortCode
	url.OriginalURL = originalURL
	url.UserID = userID
	url.Clicks = 0

	return &url, nil
}

// GetURLByShortCode retrieves a URL record by its short code
func (s *URLService) GetURLByShortCode(shortCode string) (*models.URL, error) {
	query := `
		SELECT id, short_code, original_url, user_id, clicks, created_at, updated_at
		FROM urls
		WHERE short_code = $1
	`

	var url models.URL
	err := s.db.QueryRow(query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.UserID,
		&url.Clicks,
		&url.CreatedAt,
		&url.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("URL not found")
		}

		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	return &url, nil
}

// IncrementClickCount increases the click count for a URL
func (s *URLService) IncrementClickCount(urlID int) error {
	query := `UPDATE urls SET clicks = clicks + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	_, err := s.db.Exec(query, urlID)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	return nil
}

// GetUserURLs fetches all URLs created by a user, ordered by newest first
func (s *URLService) GetUserURLs(userID int) ([]*models.URL, error) {
	query := `
		SELECT id, short_code, original_url, user_id, clicks, created_at, updated_at
		FROM urls
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user URLs: %w", err)
	}
	defer rows.Close()

	var urls []*models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.OriginalURL,
			&url.UserID,
			&url.Clicks,
			&url.CreatedAt,
			&url.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, &url)
	}

	return urls, nil
}

// generateShortCode creates a random, url safe alphanumeric short code (max length 8 chars)
func (s *URLService) generateShortCode() (string, error) {
	b := make([]rune, 8)
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[n.Int64()]
	}

	return string(b), nil
}

// shortCodeExists checks if a short code is already in the db
func (s *URLService) shortCodeExists(shortCode string) bool {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	var exists bool
	s.db.QueryRow(query, shortCode).Scan(&exists)

	return exists
}
