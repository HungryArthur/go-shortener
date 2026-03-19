package repository

import "errors"

var (
	ErrShortenedURLAlreadyExists = errors.New("shortened url already exists")
	ErrShortenedURLDoesntExist   = errors.New("shortened url doesn't exist")
	ErrCantSaveURL               = errors.New("can't save url")
	ErrCantGetSourceURL          = errors.New("can't get source url")
)
