package services

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrRateLimited       = errors.New("rate limited")
	ErrNotFound          = errors.New("not found")
)
