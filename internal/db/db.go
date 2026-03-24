package db

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

func Connect(dsn string, logger *zap.Logger) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is empty")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Создаем таблицу если не существует
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(10) UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		clicks INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	logger.Info("database connected successfully")
	return db, nil
}