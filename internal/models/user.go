package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserCredentials represents user credentials for login
type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegistration represents user registration data
type UserRegistration struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}