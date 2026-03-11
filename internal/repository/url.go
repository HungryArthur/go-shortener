package repository

import (
	"errors"
)

type UrlRepository struct {
	// сокращенный url -> куда ведет сокращенный url
	st map[string]string
}

func NewUrlRepository() *UrlRepository {
	return &UrlRepository{
		st: make(map[string]string),
	}
}

var (
	ErrShortenedUrlAlreadyExists = errors.New("shortened url already exists")
	ErrShortenedUrlDoesntExist   = errors.New("shortened url doesn't exist")
)

func (r *UrlRepository) Save(sourceUrl, shortened string) error {
	_, ok := r.st[shortened]
	if ok {
		return ErrShortenedUrlAlreadyExists
	}

	r.st[shortened] = sourceUrl
	return nil
}

func (r *UrlRepository) Get(shortenedUrl string) (string, error) {
	srcUrl, ok := r.st[shortenedUrl]
	if !ok {
		return "", ErrShortenedUrlDoesntExist
	}
	return srcUrl, nil
}
