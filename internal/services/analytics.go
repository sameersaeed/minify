package services

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"minify/internal/models"
)

type AnalyticsService struct {
	db *sql.DB
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// RecordClick logs a click asynchronously, storing user agent and ip
func (s *AnalyticsService) RecordClick(urlID int, userAgent, ipAddress string) error {
	log.Printf("[AnalyticsService] Recording click for URL ID %d\n", urlID)
	go func() {
		query := `
			INSERT INTO clicks (url_id, user_agent, ip_address)
			VALUES ($1, $2, $3)
		`
		if _, err := s.db.Exec(query, urlID, userAgent, ipAddress); err != nil {
			log.Println("[AnalyticsService] Failed to record click:", err)
		} else {
			log.Printf("[AnalyticsService] Click recorded for URL ID %d\n", urlID)
		}
	}()
	return nil
}

// GetOverview concurrently fetches the overall stats for the service
func (s *AnalyticsService) GetOverview() (*models.OverviewStats, error) {
	log.Println("[AnalyticsService] Fetching overview stats")
	var wg sync.WaitGroup
	var mu sync.Mutex
	stats := &models.OverviewStats{TimeframeData: make(map[string]interface{})}
	var err error

	wg.Add(4)

	go func() { // total users
		defer wg.Done()
		query := `SELECT COUNT(*) FROM users`
		if scanErr := s.db.QueryRow(query).Scan(&stats.TotalUsers); scanErr != nil {
			mu.Lock()
			err = fmt.Errorf("failed to get total users: %w", scanErr)
			mu.Unlock()
			log.Println("[AnalyticsService] Error fetching total users:", scanErr)
		}
	}()

	go func() { // total Minified URLs generated
		defer wg.Done()
		query := `SELECT COUNT(*) FROM urls`
		if scanErr := s.db.QueryRow(query).Scan(&stats.TotalURLs); scanErr != nil {
			mu.Lock()
			err = fmt.Errorf("failed to get total URLs: %w", scanErr)
			mu.Unlock()
			log.Println("[AnalyticsService] Error fetching total URLs:", scanErr)
		}
	}()

	go func() { // total Minified URL clicks
		defer wg.Done()
		query := `SELECT COALESCE(SUM(clicks), 0) FROM urls`
		if scanErr := s.db.QueryRow(query).Scan(&stats.TotalClicks); scanErr != nil {
			mu.Lock()
			err = fmt.Errorf("failed to get total clicks: %w", scanErr)
			mu.Unlock()
			log.Println("[AnalyticsService] Error fetching total clicks:", scanErr)
		}
	}()

	go func() { // recent users
		defer wg.Done()
		query := `SELECT username FROM users ORDER BY created_at DESC LIMIT 10`
		rows, scanErr := s.db.Query(query)
		if scanErr != nil {
			mu.Lock()
			err = fmt.Errorf("failed to get recent users: %w", scanErr)
			mu.Unlock()
			log.Println("[AnalyticsService] Error fetching recent users:", scanErr)
			return
		}
		defer rows.Close()
		var users []string
		for rows.Next() {
			var username string
			if rowErr := rows.Scan(&username); rowErr != nil {
				continue
			}
			users = append(users, username)
		}
		mu.Lock()
		stats.RecentUsers = users
		mu.Unlock()
	}()

	wg.Wait()
	if err != nil {
		return nil, err
	}

	// timeframe data
	timeframes := []string{"hour", "day", "month", "year"}
	for _, period := range timeframes {
		if data, tfErr := s.GetTimeframeStats(period); tfErr == nil {
			stats.TimeframeData[period] = data
		}
	}

	log.Println("[AnalyticsService] Overview stats fetched successfully")
	return stats, nil
}

// GetPopularURLs returns the top clicked urls, joins user info if available
func (s *AnalyticsService) GetPopularURLs(limit int) ([]*models.PopularURL, error) {
	if limit <= 0 {
		limit = 10
	}
	log.Printf("[AnalyticsService] Fetching top %d popular URLs\n", limit)

	query := `
		SELECT u.short_code, u.original_url, u.clicks, COALESCE(us.username, '') as username
		FROM urls u
		LEFT JOIN users us ON u.user_id = us.id
		ORDER BY u.clicks DESC
		LIMIT $1
	`
	rows, err := s.db.Query(query, limit)
	if err != nil {
		log.Println("[AnalyticsService] Failed to query popular URLs:", err)
		return nil, fmt.Errorf("failed to get popular URLs: %w", err)
	}
	defer rows.Close()

	var urls []*models.PopularURL
	for rows.Next() {
		var url models.PopularURL
		if scanErr := rows.Scan(&url.ShortCode, &url.OriginalURL, &url.Clicks, &url.Username); scanErr != nil {
			log.Println("[AnalyticsService] Failed to scan popular URL:", scanErr)
			continue
		}
		urls = append(urls, &url)
	}

	log.Printf("[AnalyticsService] Fetched %d popular URLs\n", len(urls))
	return urls, nil
}

// GetTimeframeStats gets clicks, urls, and unique users for a given interval
func (s *AnalyticsService) GetTimeframeStats(period string) (*models.TimeframeStats, error) {
	log.Println("[AnalyticsService] Fetching timeframe stats for period:", period)
	var interval string
	switch period {
	case "hour":
		interval = "1 hour"
	case "day":
		interval = "1 day"
	case "month":
		interval = "1 month"
	case "year":
		interval = "1 year"
	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	stats := &models.TimeframeStats{Period: period}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error

	wg.Add(3)

	go func() { // click count
		defer wg.Done()
		query := fmt.Sprintf("SELECT COUNT(*) FROM clicks WHERE clicked_at >= NOW() - INTERVAL '%s'", interval)
		if scanErr := s.db.QueryRow(query).Scan(&stats.ClickCount); scanErr != nil {
			mu.Lock()
			err = fmt.Errorf("failed to get click count: %w", scanErr)
			mu.Unlock()
			log.Println("[AnalyticsService] Error fetching click count:", scanErr)
		}
	}()

	go func() { // URL count
		defer wg.Done()
		query := fmt.Sprintf("SELECT COUNT(*) FROM urls WHERE created_at >= NOW() - INTERVAL '%s'", interval)
		if scanErr := s.db.QueryRow(query).Scan(&stats.URLCount); scanErr != nil {
			mu.Lock()
			err = fmt.Errorf("failed to get URL count: %w", scanErr)
			mu.Unlock()
			log.Println("[AnalyticsService] Error fetching URL count:", scanErr)
		}
	}()

	go func() { // unique users
		defer wg.Done()
		query := fmt.Sprintf("SELECT COUNT(DISTINCT user_id) FROM urls WHERE user_id IS NOT NULL AND created_at >= NOW() - INTERVAL '%s'", interval)
		if scanErr := s.db.QueryRow(query).Scan(&stats.UniqueUsers); scanErr != nil {
			mu.Lock()
			err = fmt.Errorf("failed to get unique users: %w", scanErr)
			mu.Unlock()
			log.Println("[AnalyticsService] Error fetching unique users:", scanErr)
		}
	}()

	wg.Wait()
	if err != nil {
		return nil, err
	}

	log.Println("[AnalyticsService] Timeframe stats fetched for period:", period)
	return stats, nil
}
