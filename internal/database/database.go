package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// Connect opens a PostgreSQL connection with reasonable defaults for
// max open/idle connections and connection lifetime, and verifies it with Ping.
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate makes sure that the required tables and indexes exist for the server to run properly. 
// If they don't then the tables are created
func Migrate(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_code VARCHAR(10) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			clicks INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS clicks (
			id SERIAL PRIMARY KEY,
			url_id INTEGER REFERENCES urls(id) ON DELETE CASCADE,
			user_agent TEXT,
			ip_address VARCHAR(45),
			clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code)`,
		`CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_clicks_url_id ON clicks(url_id)`,
		`CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}
