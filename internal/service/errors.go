package service

import "errors"

var (
	ErrShortenedURLDoesntExist  = errors.New("shortened url doesn't exist")
	ErrCantGetURL               = errors.New("can't get url")
	ErrCantSaveURL              = errors.New("can't save url")
	ErrDatabaseConnectionFailed = errors.New("database connection failed")
)
