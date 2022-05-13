package model

import "errors"

var (
	ErrNotFound     = errors.New("Item not found")
	ErrUnauthorized = errors.New("Unauthorized")
	ErrExists       = errors.New("Item already exists")
	ErrSignature    = errors.New("invalid signature")
	ErrKeyword      = errors.New("keywords do not match")
)
