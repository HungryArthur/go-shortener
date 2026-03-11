package repository

import (
	"errors"
)

type URLRepository struct {
	// сокращенный url -> куда ведет сокращенный url
	st map[string]string
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		st: make(map[string]string),
	}
}

var (
	ErrShortenedURLAlreadyExists = errors.New("shortened url already exists")
	ErrShortenedURLDoesntExist   = errors.New("shortened url doesn't exist")
)

func (r *URLRepository) Save(sourceURL, shortened string) error {
	_, ok := r.st[shortened]
	if ok {
		return ErrShortenedURLAlreadyExists
	}

	r.st[shortened] = sourceURL
	return nil
}

func (r *URLRepository) Get(shortenedURL string) (string, error) {
	srcURL, ok := r.st[shortenedURL]
	if !ok {
		return "", ErrShortenedURLDoesntExist
	}
	return srcURL, nil
}
