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

type LoginRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token       string `json:"token"`
	UserID      int64  `json:"user_id"`
	Role        string `json:"role"`
	ReferenceID int64  `json:"reference_id"`
}

type RegisterRetailerRequest struct {
	OwnerName    string `json:"name" validate:"required,min=3,max=50"`
	BusinessName string `json:"business_name" validate:"required,min=3,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required,numeric,min=10,max=15"`
	Country      string `json:"country" validate:"required"`
	State        string `json:"state" validate:"required"`
	City         string `json:"city" validate:"required"`
	PinCode      string `json:"pin_code" validate:"required"`
	Address      string `json:"address" validate:"required,min=10,max=200"`
	GSTNumber    string `json:"gst_number" validate:"required"`
	Username     string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password     string `json:"password" validate:"required,min=8"`
}

type CreateStockistUserRequest struct {
	OwnerName    string `json:"name" validate:"required,min=3,max=50"`
	BusinessName string `json:"business_name" validate:"required,min=3,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required,numeric,min=10,max=15"`
	Country      string `json:"country" validate:"required"`
	State        string `json:"state" validate:"required"`
	City         string `json:"city" validate:"required"`
	PinCode      string `json:"pin_code" validate:"required"`
	Address      string `json:"address" validate:"required,min=10,max=200"`
	GSTNumber    string `json:"gst_number" validate:"required"`
	Username     string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password     string `json:"password" validate:"required,min=8"`
}
