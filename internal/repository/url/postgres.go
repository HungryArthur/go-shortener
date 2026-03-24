package repository

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

type URLPostgresRepository struct {
	db *sql.DB
	logger *zap.Logger
}

func NewURLPostgresRepository(db *sql.DB, logger *zap.Logger) *URLPostgresRepository {
	return &URLPostgresRepository{
		db:     db,
		logger: logger.With(zap.String("repo", "postgres")),
	}
}

func (r *URLPostgresRepository) Save(sourceURL, shortenedURL string) error {
	// Проверка существует ли
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`, shortenedURL).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantSaveURL, err)
	}
	if exists {
		return ErrShortenedURLAlreadyExists
	}

	// Сохранение
	_, err = r.db.Exec(`INSERT INTO urls (short_code, original_url) VALUES ($1, $2)`, shortenedURL, sourceURL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCantSaveURL, err)
	}
	return nil
}

func (r *URLPostgresRepository) Get(shortenedURL string) (string, error) {
	var sourceURL string
	err := r.db.QueryRow(`
		UPDATE urls SET clicks = clicks + 1 
		WHERE short_code = $1 
		RETURNING original_url
	`, shortenedURL).Scan(&sourceURL)
	
	if err == sql.ErrNoRows {
		return "", ErrShortenedURLDoesntExist
	}
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCantGetSourceURL, err)
	}
	return sourceURL, nil
}

func (r *URLPostgresRepository) Ping() error {
	if r.db == nil {
		return ErrDatabaseConnectionFailed
	}
	return r.db.Ping()
}