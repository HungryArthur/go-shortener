package service

import "errors"

var (
	ErrShortenedURLDoesntExist = errors.New("shortened url doesn't exist")
	ErrCantGetURL              = errors.New("can't get url")
)
