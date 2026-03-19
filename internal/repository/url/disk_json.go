package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
)

type URLDiskJsonRepository struct {
	jsonStoragePath string
	mu              sync.Mutex

	logger *zap.Logger
}

func NewURLDiskJsonRepository(logger *zap.Logger, jsonStoragePath string) *URLDiskJsonRepository {
	return &URLDiskJsonRepository{
		jsonStoragePath: jsonStoragePath,
		mu:              sync.Mutex{},
		logger:          logger.With(zap.String("service", "url-json-repo")),
	}
}

type urlEntry struct {
	ShortenedURL string `json:"short_url"`
	SourceURL    string `json:"original_url"`
}

func (r *URLDiskJsonRepository) Save(sourceURL, shortenedURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	urlEntries := make([]urlEntry, 0)

	fileData, err := os.ReadFile(r.jsonStoragePath)

	if err != nil {
		r.logger.Error("can't read file",
			zap.String("path", r.jsonStoragePath),
			zap.String("err", err.Error()),
		)
		return fmt.Errorf("%w: can't read file: %w", err, ErrCantSaveURL)
	}

	err = json.Unmarshal(fileData, &urlEntries)
	if err != nil {
		r.logger.Error("can't unmarshal file, resetting to blank storage",
			zap.String("err", err.Error()),
		)
		urlEntries = urlEntries[:0]
	}

	for _, e := range urlEntries {
		if e.ShortenedURL == shortenedURL {
			return ErrShortenedURLAlreadyExists
		}
	}

	urlEntries = append(urlEntries, urlEntry{
		ShortenedURL: shortenedURL,
		SourceURL:    sourceURL,
	})

	data, err := json.Marshal(urlEntries)

	if err != nil {
		r.logger.Error("can't marshal url Entries",
			zap.String("err", err.Error()),
		)
		return fmt.Errorf("%w: can't marshal url data: %w", err, ErrCantSaveURL)
	}

	err = os.WriteFile(r.jsonStoragePath, data, 0644)
	if err != nil {
		r.logger.Error("can't write file",
			zap.String("err", err.Error()),
		)
		return fmt.Errorf("%w: can't save json data to file: %w", err, ErrCantSaveURL)
	}

	return nil
}

func (r *URLDiskJsonRepository) Get(shortenedURL string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	urlEntries := make([]urlEntry, 0)

	fileData, err := os.ReadFile(r.jsonStoragePath)

	if err != nil {
		r.logger.Error("can't read file",
			zap.String("path", r.jsonStoragePath),
			zap.String("err", err.Error()),
		)
		return "", fmt.Errorf("%w: can't read file: %w", err, ErrCantGetSourceURL)
	}

	err = json.Unmarshal(fileData, &urlEntries)
	if err != nil {
		r.logger.Error("can't unmarshal file, resetting it",
			zap.String("err", err.Error()),
		)
		err = os.WriteFile(r.jsonStoragePath, []byte("[]"), 0644)
		if err != nil {
			r.logger.Error("can't write file",
				zap.String("err", err.Error()),
			)
			return "", fmt.Errorf("%w: file was corrupted, but an error occurred when trying to overwrite it : %w", err, ErrCantGetSourceURL)
		}
		return "", ErrShortenedURLDoesntExist
	}
	for _, e := range urlEntries {
		if e.ShortenedURL == shortenedURL {
			return e.SourceURL, nil
		}
	}
	return "", ErrShortenedURLDoesntExist
}
