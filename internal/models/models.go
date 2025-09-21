package models

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type URL struct {
	ID          int       `json:"id" db:"id"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	UserID      *int      `json:"user_id,omitempty" db:"user_id"`
	Clicks      int       `json:"clicks" db:"clicks"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Click struct {
	ID        int       `json:"id" db:"id"`
	URLID     int       `json:"url_id" db:"url_id"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	ClickedAt time.Time `json:"clicked_at" db:"clicked_at"`
}

// requests / responses
type MinifyRequest struct {
	URL    string `json:"url" validate:"required,url"`
	UserID *int   `json:"user_id,omitempty"`
}


type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type MinifyResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// analytics 
type OverviewStats struct {
	TotalUsers    int                    `json:"total_users"`
	TotalURLs     int                    `json:"total_urls"`
	TotalClicks   int                    `json:"total_clicks"`
	RecentUsers   []string               `json:"recent_users"`
	TimeframeData map[string]interface{} `json:"timeframe_data"`
}

type PopularURL struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	Clicks      int    `json:"clicks"`
	Username    string `json:"username,omitempty"`
}

type TimeframeStats struct {
	Period      string `json:"period"`
	ClickCount  int    `json:"click_count"`
	URLCount    int    `json:"url_count"`
	UniqueUsers int    `json:"unique_users"`
}
