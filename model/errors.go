package model

import "errors"

var (
	ErrNotFound     = errors.New("Item not found")
	ErrUnauthorized = errors.New("Unauthorized")
)
