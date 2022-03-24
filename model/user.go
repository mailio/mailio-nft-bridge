package model

const UserTable = "user"

// User object (only admin for access to private methods)
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Created  int64  `json:"created,omitempty"`
	Modified int64  `json:"modified,omitempty"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// EmailPasswordInput - email and password for login purposes
type EmailPasswordInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JwtTokenOutput struct {
	Token string `json:"token"`
}
