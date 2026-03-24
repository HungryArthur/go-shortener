package repository

import "errors"

var (
	ErrShortenedURLAlreadyExists = errors.New("shortened url already exists")
	ErrCantSaveURL               = errors.New("can't save url")

	ErrShortenedURLDoesntExist   = errors.New("shortened url doesn't exist")
	ErrCantGetSourceURL          = errors.New("can't get source url")
	
	ErrDatabaseConnectionFailed = errors.New("database connection failed")
)