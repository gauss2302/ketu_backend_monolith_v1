package domain

import "errors"

// Common errors
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidID    = errors.New("invalid id")
)

// User-related errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Place-related errors
var (
	ErrPlaceNotFound     = errors.New("place not found")
	ErrPlaceExists       = errors.New("place already exists")
	ErrInvalidCriteria   = errors.New("invalid search criteria")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
	ErrEmptyName         = errors.New("name cannot be empty")
	ErrEmptyDescription  = errors.New("description cannot be empty")
	ErrEmptyAddress      = errors.New("address cannot be empty")
)
