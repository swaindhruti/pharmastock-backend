package auth

import "time"

type User struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	Role         string
	ReferenceID  *int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Claims struct {
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	ReferenceID int64  `json:"reference_id"`
}
