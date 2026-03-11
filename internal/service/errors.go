package service

import "errors"

var (
	ErrShortenedUrlDoesntExist = errors.New("shortened url doesn't exist")
	ErrCantGetUrl              = errors.New("can't get url")
)
